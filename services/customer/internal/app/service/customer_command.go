package service

import (
  "context"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  "rashop/services/customer/internal/app/command"
  "rashop/services/customer/pkg/cqrs"
)

type ICustomerCommand interface {
  AddAddress(ctx context.Context, cmd *command.AddCustomerAddressCommand) (types.Id, status.Object)
  AddVouchers(ctx context.Context, cmd *command.AddCustomerVoucherCommand) status.Object
  Create(ctx context.Context, cmd *command.CreateCustomerCommand) (types.Id, status.Object)
  DeleteAddress(ctx context.Context, cmd *command.DeleteCustomerAddressCommand) status.Object
  DeleteVouchers(ctx context.Context, cmd *command.DeleteCustomerVoucherCommand) status.Object
  Disable(ctx context.Context, cmd *command.DisableCustomerCommand) status.Object
  Enable(ctx context.Context, cmd *command.EnableCustomerCommand) status.Object
  ForgotPasswordRequest(ctx context.Context, cmd *command.ForgotCustomerPasswordRequestCommand) status.Object
  ResetPassword(ctx context.Context, cmd *command.ResetCustomerPasswordCommand) status.Object
  SetDefaultAddress(ctx context.Context, cmd *command.SetCustomerDefaultAddressCommand) status.Object
  UpdateAddress(ctx context.Context, cmd *command.UpdateCustomerAddressCommand) status.Object
  UpdateBalance(ctx context.Context, cmd *command.UpdateCustomerBalanceCommand) status.Object
  Update(ctx context.Context, cmd *command.UpdateCustomerCommand) status.Object
  UpdatePassword(ctx context.Context, cmd *command.UpdateCustomerPasswordCommand) status.Object
  UpdatePhoto(ctx context.Context, cmd *command.UpdateCustomerPhotoCommand) status.Object
  UpdateVoucher(ctx context.Context, cmd *command.UpdateCustomerVoucherCommand) status.Object
  VerifyEmailRequest(ctx context.Context, cmd *command.VerificationCustomerEmailRequestCommand) status.Object
  VerifyEmail(ctx context.Context, cmd *command.VerifyCustomerEmailCommand) status.Object
}

func NewCustomerCommand(config CustomerCommandFactory) ICustomerCommand {
  return &customerCommandService{
    i: config,
  }
}

func DefaultCustomerCommandFactory(parameter cqrs.CommonHandlerParameter) CustomerCommandFactory {
  return CustomerCommandFactory{
    AddAddress:        command.NewAddCustomerAddressHandler(parameter),
    AddVouchers:       command.NewAddCustomerVouchersHandler(parameter),
    Create:            command.NewCreateCustomerHandler(parameter),
    DeleteAddress:     command.NewDeleteCustomerAddressHandler(parameter),
    DeleteVoucher:     command.NewDeleteCustomerVoucherHandler(parameter),
    Disable:           command.NewDisableCustomerHandler(parameter),
    Enable:            command.NewEnableCustomerHandler(parameter),
    ForgotPasswordReq: command.NewForgotCustomerPasswordRequestHandler(parameter),
    ResetPassword:     command.NewResetPasswordHandler(parameter),
    SetDefaultAddress: command.NewSetCustomerDefaultAddressHandler(parameter),
    UpdateAddress:     command.NewUpdateCustomerAddressHandler(parameter),
    UpdateBalance:     command.NewUpdateCustomerBalanceHandler(parameter),
    Update:            command.NewUpdateCustomerHandler(parameter),
    UpdatePassword:    command.NewUpdatePasswordHandler(parameter),
    UpdatePhoto:       command.NewUpdatePhotoHandler(parameter),
    UpdateVoucher:     command.NewUpdateCustomerVoucherHandler(parameter),
    VerifyEmail:       command.NewVerifyCustomerEmailHandler(parameter),
    VerifyEmailReq:    command.NewVerificationCustomerEmailRequestHandler(parameter),
  }
}

type CustomerCommandFactory struct {
  AddAddress        command.IAddCustomerAddressHandler
  AddVouchers       command.IAddCustomerVouchersHandler
  Create            command.ICreateCustomerHandler
  DeleteAddress     command.IDeleteCustomerAddressHandler
  DeleteVoucher     command.IDeleteCustomerVoucherHandler
  Disable           command.IDisableCustomerHandler
  Enable            command.IEnableCustomerHandler
  ForgotPasswordReq command.IForgotCustomerPasswordRequestHandler
  ResetPassword     command.IResetPasswordHandler
  SetDefaultAddress command.ISetCustomerDefaultAddressHandler
  UpdateAddress     command.IUpdateCustomerAddressHandler
  UpdateBalance     command.IUpdateCustomerBalanceHandler
  Update            command.IUpdateCustomerHandler
  UpdatePassword    command.IUpdatePasswordHandler
  UpdatePhoto       command.IUpdatePhotoHandler
  UpdateVoucher     command.IUpdateCustomerVoucherHandler
  VerifyEmail       command.IVerifyCustomerEmailHandler
  VerifyEmailReq    command.IVerificationCustomerEmailRequestHandler
}

type customerCommandService struct {
  i CustomerCommandFactory
}

func (c *customerCommandService) AddAddress(ctx context.Context, cmd *command.AddCustomerAddressCommand) (types.Id, status.Object) {
  return c.i.AddAddress.Handle(ctx, cmd)
}

func (c *customerCommandService) AddVouchers(ctx context.Context, cmd *command.AddCustomerVoucherCommand) status.Object {
  return c.i.AddVouchers.Handle(ctx, cmd)
}

func (c *customerCommandService) DeleteAddress(ctx context.Context, cmd *command.DeleteCustomerAddressCommand) status.Object {
  return c.i.DeleteAddress.Handle(ctx, cmd)
}

func (c *customerCommandService) DeleteVouchers(ctx context.Context, cmd *command.DeleteCustomerVoucherCommand) status.Object {
  return c.i.DeleteVoucher.Handle(ctx, cmd)
}

func (c *customerCommandService) SetDefaultAddress(ctx context.Context, cmd *command.SetCustomerDefaultAddressCommand) status.Object {
  return c.i.SetDefaultAddress.Handle(ctx, cmd)
}

func (c *customerCommandService) UpdateAddress(ctx context.Context, cmd *command.UpdateCustomerAddressCommand) status.Object {
  return c.i.UpdateAddress.Handle(ctx, cmd)
}

func (c *customerCommandService) UpdateVoucher(ctx context.Context, cmd *command.UpdateCustomerVoucherCommand) status.Object {
  return c.i.UpdateVoucher.Handle(ctx, cmd)
}

func (c *customerCommandService) Create(ctx context.Context, cmd *command.CreateCustomerCommand) (types.Id, status.Object) {
  return c.i.Create.Handle(ctx, cmd)
}

func (c *customerCommandService) Disable(ctx context.Context, cmd *command.DisableCustomerCommand) status.Object {
  return c.i.Disable.Handle(ctx, cmd)
}

func (c *customerCommandService) Enable(ctx context.Context, cmd *command.EnableCustomerCommand) status.Object {
  return c.i.Enable.Handle(ctx, cmd)
}

func (c *customerCommandService) ForgotPasswordRequest(ctx context.Context, cmd *command.ForgotCustomerPasswordRequestCommand) status.Object {
  return c.i.ForgotPasswordReq.Handle(ctx, cmd)
}

func (c *customerCommandService) ResetPassword(ctx context.Context, cmd *command.ResetCustomerPasswordCommand) status.Object {
  return c.i.ResetPassword.Handle(ctx, cmd)
}

func (c *customerCommandService) Update(ctx context.Context, cmd *command.UpdateCustomerCommand) status.Object {
  return c.i.Update.Handle(ctx, cmd)
}

func (c *customerCommandService) UpdateBalance(ctx context.Context, cmd *command.UpdateCustomerBalanceCommand) status.Object {
  return c.i.UpdateBalance.Handle(ctx, cmd)
}

func (c *customerCommandService) UpdatePassword(ctx context.Context, cmd *command.UpdateCustomerPasswordCommand) status.Object {
  return c.i.UpdatePassword.Handle(ctx, cmd)
}

func (c *customerCommandService) UpdatePhoto(ctx context.Context, cmd *command.UpdateCustomerPhotoCommand) status.Object {
  return c.i.UpdatePhoto.Handle(ctx, cmd)
}

func (c *customerCommandService) VerifyEmailRequest(ctx context.Context, cmd *command.VerificationCustomerEmailRequestCommand) status.Object {
  return c.i.VerifyEmailReq.Handle(ctx, cmd)
}

func (c *customerCommandService) VerifyEmail(ctx context.Context, cmd *command.VerifyCustomerEmailCommand) status.Object {
  return c.i.VerifyEmail.Handle(ctx, cmd)
}
