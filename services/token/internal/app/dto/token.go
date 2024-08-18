package dto

import (
  vob "github.com/arcorium/rashop/services/token/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type TokenResponseDTO struct {
  Id        types.Id
  Token     string
  Usage     vob.TokenUsage
  Type      vob.TokenType
  ExpiredAt time.Time
}
