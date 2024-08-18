package mapper

import (
  tokenv1 "github.com/arcorium/rashop/proto/gen/go/token/v1"
  "github.com/arcorium/rashop/services/token/internal/app/command"
  "github.com/arcorium/rashop/shared/errors"
  "github.com/arcorium/rashop/shared/types"
)

func ToGenerateTokenCommand(request *tokenv1.GenerateTokenRequest) (command.GenerateTokenCommand, error) {
  var fieldErrors []errors.FieldError

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    fieldErrors = append(fieldErrors, errors.NewFieldError("user_id", err))
  }

  tokenType, err := toDomainType(request.Type)
  if err != nil {
    fieldErrors = append(fieldErrors, errors.NewFieldError("type", err))
  }

  tokenUsage, err := toDomainUsage(request.Usage)
  if err != nil {
    fieldErrors = append(fieldErrors, errors.NewFieldError("usage", err))
  }

  if len(fieldErrors) > 0 {
    return command.GenerateTokenCommand{}, errors.GrpcFieldErrors(fieldErrors...)
  }

  var length uint16 = 0
  if request.Length != nil {
    length = uint16(*request.Length)
  }

  return command.GenerateTokenCommand{
    UserId:      userId,
    TokenLength: length,
    Type:        tokenType,
    Usage:       tokenUsage,
  }, nil
}

func ToVerifyTokenCommand(request *tokenv1.ValidateTokenRequest) (command.VerifyTokenCommand, error) {
  tokenUsage, err := toDomainUsage(request.Usage)
  if err != nil {
    return command.VerifyTokenCommand{}, errors.NewFieldError("usage", err).ToGrpcError()
  }

  return command.VerifyTokenCommand{
    Token: request.Token,
    Usage: tokenUsage,
  }, nil
}
