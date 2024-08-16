package mapper

import (
  customerv1 "github.com/arcorium/rashop/proto/gen/go/customer/v1"
  sharedErr "github.com/arcorium/rashop/shared/errors"
  "github.com/arcorium/rashop/shared/types"
  "github.com/arcorium/rashop/shared/util"
  "rashop/services/customer/internal/app/command"
)

func ToAddCustomerAddressCommand(request *customerv1.AddCustomerAddressRequest) (command.AddCustomerAddressCommand, error) {
  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    return command.AddCustomerAddressCommand{}, sharedErr.NewFieldError("customer_id", err).ToGrpcError()
  }

  res := command.AddCustomerAddressCommand{
    CustomerId:     custId,
    StreetAddress1: request.StreetAddress_1,
    StreetAddress2: types.NewNullable(request.StreetAddress_2),
    City:           request.City,
    State:          request.State,
    PostalCode:     request.PostalCode,
  }

  err = util.ValidateStruct(&res)
  return res, err
}

func ToAddCustomerVoucherCommand(request *customerv1.AddCustomerVoucherRequest) (command.AddCustomerVoucherCommand, error) {
  var fieldsErr []sharedErr.FieldError

  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("customer_id", err))
  }

  voucherId, err := types.IdFromString(request.VoucherId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("voucher_id", err))
  }

  if len(fieldsErr) > 0 {
    return command.AddCustomerVoucherCommand{}, sharedErr.GrpcFieldErrors(fieldsErr...)
  }

  return command.AddCustomerVoucherCommand{
    CustomerId: custId,
    VoucherId:  voucherId,
  }, nil
}

func ToDeleteCustomerAddressCommand(request *customerv1.DeleteCustomerAddressRequest) (command.DeleteCustomerAddressCommand, error) {
  var fieldsErr []sharedErr.FieldError

  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("customer_id", err))
  }

  addressId, err := types.IdFromString(request.AddressId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("address_id", err))
  }

  if len(fieldsErr) > 0 {
    return command.DeleteCustomerAddressCommand{}, sharedErr.GrpcFieldErrors(fieldsErr...)
  }

  return command.DeleteCustomerAddressCommand{
    CustomerId: custId,
    AddressId:  addressId,
  }, nil
}

func ToDeleteCustomerVoucherCommand(request *customerv1.DeleteCustomerVoucherRequest) (command.DeleteCustomerVoucherCommand, error) {
  var fieldsErr []sharedErr.FieldError

  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("customer_id", err))
  }

  addressId, err := types.IdFromString(request.VoucherId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("voucher_id", err))
  }

  if len(fieldsErr) > 0 {
    return command.DeleteCustomerVoucherCommand{}, sharedErr.GrpcFieldErrors(fieldsErr...)
  }

  return command.DeleteCustomerVoucherCommand{
    CustomerId: custId,
    VoucherId:  addressId,
  }, nil
}

func ToForgotCustomerPasswordRequestCommand(request *customerv1.ForgotCustomerPasswordRequest) (command.ForgotCustomerPasswordRequestCommand, error) {
  email, err := types.EmailFromString(request.Email)
  if err != nil {
    return command.ForgotCustomerPasswordRequestCommand{}, sharedErr.NewFieldError("email", err).ToGrpcError()
  }

  return command.ForgotCustomerPasswordRequestCommand{
    Email: email,
  }, nil
}

func ToResetCustomerPasswordCommand(request *customerv1.ResetCustomerPasswordRequest) command.ResetCustomerPasswordCommand {
  return command.ResetCustomerPasswordCommand{
    Token:       request.Token,
    LogoutAll:   request.LogoutAll,
    NewPassword: types.PasswordFromString(request.NewPassword),
  }
}

func ToSetCustomerDefaultAddressCommand(request *customerv1.SetCustomerDefaultAddressRequest) (command.SetCustomerDefaultAddressCommand, error) {
  var fieldsErr []sharedErr.FieldError

  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("customer_id", err))
  }

  addressId, err := types.IdFromString(request.AddressId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("address_id", err))
  }

  if len(fieldsErr) > 0 {
    return command.SetCustomerDefaultAddressCommand{}, sharedErr.GrpcFieldErrors(fieldsErr...)
  }

  return command.SetCustomerDefaultAddressCommand{
    CustomerId: custId,
    AddressId:  addressId,
  }, nil
}

func ToUpdateCustomerAddressCommand(request *customerv1.UpdateCustomerAddressRequest) (command.UpdateCustomerAddressCommand, error) {
  var fieldsErr []sharedErr.FieldError

  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("customer_id", err))
  }

  addressId, err := types.IdFromString(request.AddressId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("address_id", err))
  }

  if len(fieldsErr) > 0 {
    return command.UpdateCustomerAddressCommand{}, sharedErr.GrpcFieldErrors(fieldsErr...)
  }

  return command.UpdateCustomerAddressCommand{
    CustomerId:     custId,
    AddressId:      addressId,
    StreetAddress1: types.NewNullable(request.StreetAddress_1),
    StreetAddress2: types.NewNullable(request.StreetAddress_2),
    City:           types.NewNullable(request.City),
    State:          types.NewNullable(request.State),
    PostalCode:     types.NewNullable(request.PostalCode),
  }, nil
}

func ToSetCustomerBalanceCommand(request *customerv1.SetCustomerBalanceRequest) (command.UpdateCustomerBalanceCommand, error) {
  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    return command.UpdateCustomerBalanceCommand{}, sharedErr.NewFieldError("customer_id", err).ToGrpcError()
  }

  return command.UpdateCustomerBalanceCommand{
    CustomerId: custId,
    Operator:   command.OperatorSet,
    Balance:    int64(request.Balance),
    Point:      int64(request.Point),
  }, nil
}

func ToModifyCustomerBalanceCommand(request *customerv1.ModifyCustomerBalanceRequest) (command.UpdateCustomerBalanceCommand, error) {
  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    return command.UpdateCustomerBalanceCommand{}, sharedErr.NewFieldError("customer_id", err).ToGrpcError()
  }

  return command.UpdateCustomerBalanceCommand{
    CustomerId: custId,
    Operator:   command.OperatorMod,
    Balance:    request.Balance,
    Point:      request.Point,
  }, nil
}

func ToUpdateCustomerPasswordCommand(request *customerv1.UpdateCustomerPasswordRequest) (command.UpdateCustomerPasswordCommand, error) {
  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    return command.UpdateCustomerPasswordCommand{}, sharedErr.NewFieldError("customer_id", err).ToGrpcError()
  }

  return command.UpdateCustomerPasswordCommand{
    CustomerId:   custId,
    LastPassword: types.PasswordFromString(request.LastPassword),
    NewPassword:  types.PasswordFromString(request.NewPassword),
  }, nil
}

func ToUpdateCustomerPhotoCommand(request *customerv1.UpdateCustomerPhotoRequest) (command.UpdateCustomerPhotoCommand, error) {
  var fieldsErr []sharedErr.FieldError

  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("customer_id", err))
  }

  mediaId, err := types.IdFromString(request.MediaId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("media_id", err))
  }

  if len(fieldsErr) > 0 {
    return command.UpdateCustomerPhotoCommand{}, sharedErr.GrpcFieldErrors(fieldsErr...)
  }

  return command.UpdateCustomerPhotoCommand{
    CustomerId: custId,
    MediaId:    mediaId,
  }, nil
}

func ToUpdateCustomerVoucherCommand(request *customerv1.UpdateCustomerVoucherRequest) (command.UpdateCustomerVoucherCommand, error) {
  var fieldsErr []sharedErr.FieldError

  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("customer_id", err))
  }

  voucherId, err := types.IdFromString(request.VoucherId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("media_id", err))
  }

  if len(fieldsErr) > 0 {
    return command.UpdateCustomerVoucherCommand{}, sharedErr.GrpcFieldErrors(fieldsErr...)
  }

  return command.UpdateCustomerVoucherCommand{
    CustomerId:  custId,
    VoucherId:   voucherId,
    IsBeingUsed: false,
  }, nil
}

func ToVerifyCustomerEmailCommand(request *customerv1.VerifyCustomerEmailRequest) (command.VerifyCustomerEmailCommand, error) {
  return command.VerifyCustomerEmailCommand{
    Token: request.Token,
  }, nil
}

func ToVerifyCustomerEmailRequestCommand(request *customerv1.VerifyCustomerEmailInstantiateRequest) (command.VerificationCustomerEmailRequestCommand, error) {
  custId, err := types.IdFromString(request.CustomerId)
  if err != nil {
    return command.VerificationCustomerEmailRequestCommand{}, sharedErr.NewFieldError("customer_id", err)
  }

  return command.VerificationCustomerEmailRequestCommand{
    CustomerId: custId,
  }, nil
}

func ToCreateCustomerCommand(req *customerv1.CreateCustomerRequest) (command.CreateCustomerCommand, error) {
  email, err := types.EmailFromString(req.Email)
  if err != nil {
    return command.CreateCustomerCommand{}, sharedErr.NewFieldError("email", err).ToGrpcError()
  }

  return command.CreateCustomerCommand{
    Username:  req.Username,
    FirstName: req.FirstName,
    LastName:  types.NewNullable(req.LastName),
    Email:     email,
    Password:  types.PasswordFromString(req.Password),
  }, nil
}

func ToUpdateCustomerCommand(req *customerv1.UpdateCustomerRequest) (command.UpdateCustomerCommand, error) {
  var fieldsErr []sharedErr.FieldError

  custId, err := types.IdFromString(req.CustomerId)
  if err != nil {
    fieldsErr = append(fieldsErr, sharedErr.NewFieldError("customer_id", err))
  }

  var email *types.Email
  if req.Email != nil {
    emails, err := types.EmailFromString(*req.Email)
    if err != nil {
      fieldsErr = append(fieldsErr, sharedErr.NewFieldError("email", err))
    }
    email = &emails
  }

  if len(fieldsErr) > 0 {
    return command.UpdateCustomerCommand{}, sharedErr.GrpcFieldErrors(fieldsErr...)
  }

  return command.UpdateCustomerCommand{
    CustomerId: custId,
    Username:   types.NewNullable(req.Username),
    FirstName:  types.NewNullable(req.FirstName),
    LastName:   types.NewNullable(req.LastName),
    Email:      types.NewNullable(email),
  }, nil
}

func ToDisableCustomerCommand(req *customerv1.DisableCustomerRequest) (command.DisableCustomerCommand, error) {
  custId, err := types.IdFromString(req.CustomerId)
  if err != nil {
    return command.DisableCustomerCommand{}, sharedErr.NewFieldError("customer_id", err).ToGrpcError()
  }

  return command.DisableCustomerCommand{
    CustomerId: custId,
  }, nil
}

func ToEnableCustomerCommand(req *customerv1.EnableCustomerRequest) (command.EnableCustomerCommand, error) {
  custId, err := types.IdFromString(req.CustomerId)
  if err != nil {
    return command.EnableCustomerCommand{}, sharedErr.NewFieldError("customer_id", err).ToGrpcError()
  }

  return command.EnableCustomerCommand{
    CustomerId: custId,
  }, nil
}
