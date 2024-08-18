package mapper

import (
  tokenv1 "github.com/arcorium/rashop/proto/gen/go/token/v1"
  "github.com/arcorium/rashop/services/token/internal/app/dto"
  vob "github.com/arcorium/rashop/services/token/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/errors"
  "google.golang.org/protobuf/types/known/timestamppb"
)

func toDomainUsage(usage tokenv1.TokenUsage) (vob.TokenUsage, error) {
  switch usage {
  case tokenv1.TokenUsage_EmailVerification:
    return vob.UsageEmailVerification, nil
  case tokenv1.TokenUsage_ResetPassword:
    return vob.UsageResetPassword, nil
  case tokenv1.TokenUsage_Login:
    return vob.UsageLogin, nil
  case tokenv1.TokenUsage_General:
    return vob.UsageGeneral, nil
  }
  return 0, errors.ErrEnumOutOfBounds
}

func toDomainType(types tokenv1.TokenType) (vob.TokenType, error) {
  switch types {
  case tokenv1.TokenType_String:
    return vob.TypeString, nil
  case tokenv1.TokenType_PIN:
    return vob.TypeString, nil
  case tokenv1.TokenType_AlphanumericPIN:
    return vob.TypeString, nil
  }
  return 0, errors.ErrEnumOutOfBounds
}

func toProtoTokenUsage(usage vob.TokenUsage) tokenv1.TokenUsage {
  switch usage {
  case vob.UsageEmailVerification:
    return tokenv1.TokenUsage_EmailVerification
  case vob.UsageResetPassword:
    return tokenv1.TokenUsage_ResetPassword
  case vob.UsageLogin:
    return tokenv1.TokenUsage_Login
  case vob.UsageGeneral:
    return tokenv1.TokenUsage_General
  }
  panic("unreachable")
}

func toProtoTokenType(types vob.TokenType) tokenv1.TokenType {
  switch types {
  case vob.TypeString:
    return tokenv1.TokenType_String
  case vob.TypePIN:
    return tokenv1.TokenType_PIN
  case vob.TypeAlphanumericPIN:
    return tokenv1.TokenType_AlphanumericPIN
  }
  panic("unreachable")
}

func ToProtoToken(dto *dto.TokenResponseDTO) *tokenv1.Token {
  return &tokenv1.Token{
    Id:        dto.Id.String(),
    Token:     dto.Token,
    Usage:     toProtoTokenUsage(dto.Usage),
    Type:      toProtoTokenType(dto.Type),
    ExpiredAt: timestamppb.New(dto.ExpiredAt),
  }
}
