package command

import (
  "github.com/arcorium/rashop/services/token/internal/domain/entity"
  vob "github.com/arcorium/rashop/services/token/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type GenerateTokenCommand struct {
  UserId      types.Id
  TokenLength uint16
  Type        vob.TokenType
  Usage       vob.TokenUsage
}

func (g *GenerateTokenCommand) ToDomain(id types.Id, token string, expiryDuration time.Duration) entity.Token {
  currTime := time.Now()
  return entity.Token{
    Id:          id,
    Token:       token,
    UserId:      g.UserId,
    Usage:       g.Usage,
    Type:        g.Type,
    ExpiredAt:   currTime.Add(expiryDuration),
    GeneratedAt: currTime,
  }
}

type VerifyTokenCommand struct {
  Token string
  Usage vob.TokenUsage
}
