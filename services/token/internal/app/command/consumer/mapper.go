package consumer

import (
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/token/internal/app/command"
  vob "github.com/arcorium/rashop/services/token/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
)

func MapCustomerEmailVerificationRequested(ev *intev.CustomerEmailVerificationRequestedV1) (command.GenerateTokenCommand, error) {
  userId, err := types.IdFromString(ev.CustomerId)
  if err != nil {
    return command.GenerateTokenCommand{}, err
  }

  return command.GenerateTokenCommand{
    UserId:      userId,
    TokenLength: 0,
    Type:        vob.TypeString,
    Usage:       vob.UsageEmailVerification,
  }, nil
}

func MapResetPasswordRequested(ev *intev.CustomerResetPasswordRequestedV1) (command.GenerateTokenCommand, error) {
  userId, err := types.IdFromString(ev.CustomerId)
  if err != nil {
    return command.GenerateTokenCommand{}, err
  }

  return command.GenerateTokenCommand{
    UserId:      userId,
    TokenLength: 0,
    Type:        vob.TypeString,
    Usage:       vob.UsageResetPassword,
  }, nil
}

func MapLoginTokenRequested(ev *intev.LoginTokenRequestedV1) (command.GenerateTokenCommand, error) {
  userId, err := types.IdFromString(ev.UserId)
  if err != nil {
    return command.GenerateTokenCommand{}, err
  }

  return command.GenerateTokenCommand{
    UserId:      userId,
    TokenLength: 0,
    Type:        vob.TypeAlphanumericPIN,
    Usage:       vob.UsageLogin,
  }, nil
}

func MapCustomerCreated(ev *intev.CustomerCreatedV1) (command.GenerateTokenCommand, error) {
  userId, err := types.IdFromString(ev.CustomerId)
  if err != nil {
    return command.GenerateTokenCommand{}, err
  }

  return command.GenerateTokenCommand{
    UserId:      userId,
    TokenLength: 0,
    Type:        vob.TypeString,
    Usage:       vob.UsageEmailVerification,
  }, nil
}

func MapCustomerEmailChanged(ev *intev.CustomerEmailUpdatedV1) (command.GenerateTokenCommand, error) {
  userId, err := types.IdFromString(ev.CustomerId)
  if err != nil {
    return command.GenerateTokenCommand{}, err
  }

  return command.GenerateTokenCommand{
    UserId:      userId,
    TokenLength: 0,
    Type:        vob.TypeString,
    Usage:       vob.UsageEmailVerification,
  }, nil
}
