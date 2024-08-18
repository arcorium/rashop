package persistence

import (
  "context"
  "database/sql"
  "errors"
  "fmt"
  "github.com/arcorium/rashop/services/token/internal/domain/entity"
  "github.com/arcorium/rashop/services/token/internal/domain/repository"
  "github.com/arcorium/rashop/services/token/internal/infra/model"
  "github.com/arcorium/rashop/services/token/pkg/tracer"
  "github.com/arcorium/rashop/shared/types"
  "github.com/arcorium/rashop/shared/util/repo"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "github.com/redis/go-redis/v9"
  "go.opentelemetry.io/otel/trace"
  "strings"
  "time"
)

var defaultConfig = RedisTokenConfig{
  TokenNamespace: "token",
  UserNamespace:  "user",
  MaxRetries:     5,
}

func NewRedisToken(client *redis.Client, config *RedisTokenConfig) repository.ITokenPersistent {
  return &redisTokenPersistent{
    RedisTokenConfig: types.OnNil(config, defaultConfig),
    client:           client,
    tracer:           tracer.Get(),
  }
}

type RedisTokenConfig struct {
  TokenNamespace string
  UserNamespace  string
  MaxRetries     int
}

type redisTokenPersistent struct {
  RedisTokenConfig
  client *redis.Client
  tracer trace.Tracer
}

func (r *redisTokenPersistent) tokenKey(token string) string {
  return fmt.Sprintf("%s:%s", r.TokenNamespace, token)
}

func (r *redisTokenPersistent) userKey(id string) string {
  return fmt.Sprintf("%s:%s", r.UserNamespace, id)
}

// watchRetry Run command Watch and retry when fails
func (r *redisTokenPersistent) watchRetry(ctx context.Context, f func(tx *redis.Tx) error, keys ...string) error {
  var err error
  for i := 0; i < r.MaxRetries; i++ {
    err = r.client.Watch(ctx, f, keys...)
    if err == nil {
      break
    }
    // Failed due to the key is changed
    if errors.Is(err, redis.TxFailedErr) {
      continue
    }
    return err
  }
  return nil
}

func (r *redisTokenPersistent) checkCmdsWithSpan(span trace.Span, cmds ...redis.Cmder) error {
  for _, cmd := range cmds {
    if err := cmd.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return err
    }
  }
  return nil
}

func (r *redisTokenPersistent) checkCmds(cmds ...redis.Cmder) error {
  for _, cmd := range cmds {
    if err := cmd.Err(); err != nil {
      return err
    }
  }
  return nil
}

func (r *redisTokenPersistent) scanModel(token string, cmd *redis.MapStringStringCmd, durationCmd *redis.DurationCmd) (model.Token, error) {
  var models model.Token
  err := cmd.Scan(&models)
  if err != nil {
    return model.Token{}, err
  }

  times := durationCmd.Val()
  models.ExpiredAt = time.Unix(int64(times.Seconds()), 0)
  models.Token = token
  return models, nil
}

func (r *redisTokenPersistent) FindByToken(ctx context.Context, token string) (entity.Token, error) {
  ctx, span := r.tracer.Start(ctx, "redisTokenPersistent.FindByToken")
  defer span.End()

  tokenKey := r.tokenKey(token)

  var hResult *redis.MapStringStringCmd
  var durResult *redis.DurationCmd

  tx := func(tx *redis.Tx) error {
    // Check if token exists
    result, err := r.client.Exists(ctx, tokenKey).Result()
    if err != nil {
      return err
    }
    if result == 0 {
      return sql.ErrNoRows
    }

    // Get the record
    cmders, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
      hResult = pipe.HGetAll(ctx, tokenKey)
      durResult = pipe.ExpireTime(ctx, tokenKey)
      return nil
    })
    if err != nil {
      return err
    }
    return r.checkCmds(cmders...)
  }

  if err := r.watchRetry(ctx, tx, tokenKey); err != nil {
    spanUtil.RecordError(err, span)
    return entity.Token{}, err
  }

  // Scan model
  models, err := r.scanModel(token, hResult, durResult)
  if err != nil {
    spanUtil.RecordError(err, span)
    return entity.Token{}, err
  }

  return models.ToDomain()
}

func (r *redisTokenPersistent) FindByUserId(ctx context.Context, userId types.Id) ([]entity.Token, error) {
  ctx, span := r.tracer.Start(ctx, "redisTokenPersistent.FindByUserId")
  defer span.End()

  userKey := r.userKey(userId.String())
  result := make(map[string]types.Pair[*redis.MapStringStringCmd, *redis.DurationCmd])

  tx := func(tx *redis.Tx) error {
    // Scan all token key on user namespace
    var cursor uint64 = 0
    var tokenKeys []string
    for {
      // Scan
      keys, newCursor, err := tx.SScan(ctx, userKey, cursor, "*", 0).Result()
      if err != nil {
        return err
      }

      tokenKeys = append(tokenKeys, keys...)
      if cursor == 0 {
        break
      }
      cursor = newCursor
    }

    if len(tokenKeys) == 0 {
      return sql.ErrNoRows
    }

    cmders, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
      for _, tokenKey := range tokenKeys {
        // Get object
        cmd := pipe.HGetAll(ctx, tokenKey)
        // Get expiration time
        expiry := pipe.ExpireTime(ctx, tokenKey)
        result[tokenKey] = types.NewPair(cmd, expiry)
      }
      return nil
    })
    if err != nil {
      return err
    }
    return r.checkCmds(cmders...)
  }

  if err := r.watchRetry(ctx, tx, userKey); err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  // Scan result
  var entities []entity.Token
  for key, result := range result {
    // Hash result
    if err := result.First.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return nil, err
    }
    // Duration

    if err := result.Second.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return nil, err
    }

    // Token
    split := strings.Split(key, ":")
    if len(split) != 2 {
      err := errors.New("invalid token key")
      spanUtil.RecordError(err, span)
      return nil, err
    }

    // Parse token
    token, err := r.scanModel(split[1], result.First, result.Second)
    if err != nil {
      spanUtil.RecordError(err, span)
      return nil, err
    }

    domain, err := token.ToDomain()
    if err != nil {
      spanUtil.RecordError(err, span)
      return nil, err
    }

    entities = append(entities, domain)
  }
  return entities, nil
}

func (r *redisTokenPersistent) Create(ctx context.Context, token *entity.Token) error {
  ctx, span := r.tracer.Start(ctx, "redisTokenPersistent.Create")
  defer span.End()

  if !token.Created() {
    return nil
  }

  models := model.FromDomainToken(token)
  tokenKey := r.tokenKey(token.Token)
  userKey := r.userKey(models.UserId)

  // Prevent update
  count, err := r.client.Exists(ctx, tokenKey).Result()
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  if count != 0 {
    err = repo.ErrAlreadyExists
    spanUtil.RecordError(err, span)
    return err
  }

  cmds, err := r.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
    // Create token
    p.HSet(ctx, tokenKey, &models)
    // Add it into user namespace
    p.SAdd(ctx, userKey, tokenKey)
    // Set expiration time
    p.ExpireAt(ctx, tokenKey, models.ExpiredAt)
    return nil
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return r.checkCmdsWithSpan(span, cmds...)
}

func (r *redisTokenPersistent) Delete(ctx context.Context, token *entity.Token) error {
  ctx, span := r.tracer.Start(ctx, "redisTokenPersistent.Delete")
  defer span.End()

  if !token.Deleted() {
    return nil
  }

  userKey := r.userKey(token.UserId.String())
  tokenKey := r.tokenKey(token.Token)

  tx := func(tx *redis.Tx) error {
    isMember, err := tx.SIsMember(ctx, userKey, tokenKey).Result()
    if err != nil {
      return err
    }
    if !isMember {
      return sql.ErrNoRows
    }

    result, err := tx.SCard(ctx, userKey).Result()
    if err != nil {
      return err
    }

    cmds, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
      // Delete token key from user namespace
      pipe.SRem(ctx, userKey, tokenKey)
      // Delete token key
      pipe.Del(ctx, tokenKey)

      // Delete user key when there are no token key on it
      if result-1 <= 0 {
        pipe.Del(ctx, userKey)
      }
      return nil
    })

    if err != nil {
      return err
    }
    return r.checkCmdsWithSpan(span, cmds...)
  }

  if err := r.watchRetry(ctx, tx, userKey); err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}
