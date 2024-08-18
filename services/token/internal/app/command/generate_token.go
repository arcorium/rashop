package command

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  . "github.com/arcorium/rashop/services/token/internal/app/common"
  "github.com/arcorium/rashop/services/token/internal/app/dto"
  "github.com/arcorium/rashop/services/token/internal/app/mapper"
  "github.com/arcorium/rashop/services/token/internal/domain/entity"
  vob "github.com/arcorium/rashop/services/token/internal/domain/valueobject"
  sharedErr "github.com/arcorium/rashop/shared/errors"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "time"
)

type IGenerateTokenHandler interface {
  handler.Command[*GenerateTokenCommand, dto.TokenResponseDTO]
}

func NewGenerateTokenHandler(parameter CommonHandlerParameter, config TokenGenerationConfig) IGenerateTokenHandler {
  return &generateTokenHandler{
    config:        config,
    commonHandler: newCommonHandler(&parameter),
  }
}

func NewSingleTokenGenerationConfig(duration time.Duration) TokenGenerationConfig {
  return TokenGenerationConfig{
    VerificationTokenExpiration: duration,
    ResetTokenExpiration:        duration,
    LoginTokenExpiration:        duration,
    GeneralTokenExpiration:      duration,
  }
}

type TokenGenerationConfig struct {
  VerificationTokenExpiration time.Duration
  ResetTokenExpiration        time.Duration
  LoginTokenExpiration        time.Duration
  GeneralTokenExpiration      time.Duration
}

type generateTokenHandler struct {
  config TokenGenerationConfig
  commonHandler
}

func (g *generateTokenHandler) getExpiryDuration(usage vob.TokenUsage) (time.Duration, error) {
  switch usage {
  case vob.UsageEmailVerification:
    return g.config.VerificationTokenExpiration, nil
  case vob.UsageResetPassword:
    return g.config.ResetTokenExpiration, nil
  case vob.UsageLogin:
    return g.config.LoginTokenExpiration, nil
  case vob.UsageGeneral:
    return g.config.GeneralTokenExpiration, nil
  default:
    return 0, sharedErr.ErrEnumOutOfBounds
  }
}

func (g *generateTokenHandler) Handle(ctx context.Context, cmd *GenerateTokenCommand) (dto.TokenResponseDTO, status.Object) {
  ctx, span := g.tracer.Start(ctx, "generateTokenHandler.Handle")
  defer span.End()

  // Generate token
  tokens := GenerateToken(cmd.Type, cmd.TokenLength)
  expiryDuration, err := g.getExpiryDuration(cmd.Usage)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.ErrBadRequest(err)
  }

  id, err := types.NewId()
  if err != nil {
    return dto.TokenResponseDTO{}, status.ErrInternal(err)
  }

  token := cmd.ToDomain(id, tokens, expiryDuration)
  token, err = entity.CreateToken(&token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.ErrInternal(err)
  }

  // Persistent
  err = g.persistent.Create(ctx, &token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.FromRepository(err)
  }

  // Integration Event
  ev := intev.NewTokenCreated(token.UserId,
    token.Id,
    token.Token,
    token.Usage.Underlying(),
    token.Type.Underlying(),
    token.ExpiredAt,
    token.GeneratedAt)
  token.AddEvents(ev)

  // Publish
  err = g.publisher.Publish(ctx, &token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.ErrInternal(err)
  }

  return mapper.ToTokenResponseDTO(&token), status.Created()
}
