package mapper

import (
  "github.com/arcorium/rashop/services/token/internal/app/dto"
  "github.com/arcorium/rashop/services/token/internal/domain/entity"
)

func ToTokenResponseDTO(token *entity.Token) dto.TokenResponseDTO {
  return dto.TokenResponseDTO{
    Id:        token.Id,
    Token:     token.Token,
    Usage:     token.Usage,
    Type:      token.Type,
    ExpiredAt: token.ExpiredAt,
  }
}
