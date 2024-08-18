package model

import (
  "github.com/arcorium/rashop/services/token/internal/domain/entity"
  vob "github.com/arcorium/rashop/services/token/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  "github.com/arcorium/rashop/shared/util/repo"
  "time"
)

type MetadataMapOption = repo.DataAccessModelMapOption[*entity.Token, *Token]

func FromDomainToken(token *entity.Token, opts ...MetadataMapOption) Token {
  result := Token{
    Token:     token.Token,
    UserId:    token.UserId.String(),
    Usage:     token.Usage.Underlying(),
    Type:      token.Type.Underlying(),
    ExpiredAt: token.ExpiredAt,
    CreatedAt: token.GeneratedAt,
  }

  for _, opt := range opts {
    opt(token, &result)
  }
  return result
}

type Token struct {
  Token     string    `redis:"token"`
  UserId    string    `redis:"user_id"`
  Usage     uint8     `redis:"usage"`
  Type      uint8     `redis:"type"`
  ExpiredAt time.Time `redis:"-"`
  CreatedAt time.Time `redis:"created_at"`
}

func (t *Token) ToDomain() (entity.Token, error) {
  userId, err := types.IdFromString(t.UserId)
  if err != nil {
    return entity.Token{}, err
  }

  usage, err := vob.NewUsage(t.Usage)
  if err != nil {
    return entity.Token{}, err
  }

  tokenType, err := vob.NewType(t.Type)
  if err != nil {
    return entity.Token{}, err
  }

  return entity.Token{
    Token:       t.Token,
    UserId:      userId,
    Usage:       usage,
    Type:        tokenType,
    ExpiredAt:   t.ExpiredAt,
    GeneratedAt: t.CreatedAt,
  }, nil
}
