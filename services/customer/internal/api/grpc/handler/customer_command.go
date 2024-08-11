package handler

import (
  "context"
  "github.com/arcorium/rashop/proto/gen/go/customer/v1"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "mini-shop/services/user/internal/api/grpc/mapper"
  "mini-shop/services/user/internal/app/service"
  "mini-shop/services/user/pkg/tracer"
)

func NewCustomerCommand(svc service.ICustomerCommand) CustomerCommandHandler {
  return CustomerCommandHandler{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type CustomerCommandHandler struct {
  customerv1.UnimplementedCustomerCommandServiceServer

  svc    service.ICustomerCommand
  tracer trace.Tracer
}

func (c *CustomerCommandHandler) Register(server *grpc.Server) {
  customerv1.RegisterCustomerCommandServiceServer(server, c)
}

func (c *CustomerCommandHandler) AddAddress(ctx context.Context, request *customerv1.AddCustomerAddressRequest) (*customerv1.AddCustomerAddressResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.AddAddress")
  defer span.End()

  cmd, err := mapper.ToAddCustomerAddressCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, stat := c.svc.AddAddress(ctx, &cmd)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &customerv1.AddCustomerAddressResponse{
    AddressId: result.String(),
  }
  return resp, nil
}

func (c *CustomerCommandHandler) AddVoucher(ctx context.Context, request *customerv1.AddCustomerVoucherRequest) (*customerv1.AddCustomerVoucherResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.AddVoucher")
  defer span.End()

  cmd, err := mapper.ToAddCustomerVoucherCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.AddVouchers(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) DeleteAddress(ctx context.Context, request *customerv1.DeleteCustomerAddressRequest) (*customerv1.DeleteCustomerAddressResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.DeleteAddress")
  defer span.End()

  cmd, err := mapper.ToDeleteCustomerAddressCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.DeleteAddress(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) DeleteVoucher(ctx context.Context, request *customerv1.DeleteCustomerVoucherRequest) (*customerv1.DeleteCustomerVoucherResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.DeleteVoucher")
  defer span.End()

  cmd, err := mapper.ToDeleteCustomerVoucherCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.DeleteVouchers(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) ForgotPasswordInstantiate(ctx context.Context, request *customerv1.ForgotCustomerPasswordRequest) (*customerv1.ForgotCustomerPasswordResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.ForgotPasswordInstantiate")
  defer span.End()

  cmd, err := mapper.ToForgotCustomerPasswordRequestCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.ForgotPasswordRequest(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) ResetPassword(ctx context.Context, request *customerv1.ResetCustomerPasswordRequest) (*customerv1.ResetCustomerPasswordResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.ResetPassword")
  defer span.End()

  cmd := mapper.ToResetCustomerPasswordCommand(request)

  stat := c.svc.ResetPassword(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) SetDefaultAddress(ctx context.Context, request *customerv1.SetCustomerDefaultAddressRequest) (*customerv1.SetCustomerDefaultAddressResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.SetDefaultAddress")
  defer span.End()

  cmd, err := mapper.ToSetCustomerDefaultAddressCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.SetDefaultAddress(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) UpdateAddress(ctx context.Context, request *customerv1.UpdateCustomerAddressRequest) (*customerv1.UpdateCustomerAddressResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.UpdateAddress")
  defer span.End()

  cmd, err := mapper.ToUpdateCustomerAddressCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.UpdateAddress(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) SetBalance(ctx context.Context, request *customerv1.SetCustomerBalanceRequest) (*customerv1.SetCustomerBalanceResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.SetBalance")
  defer span.End()

  cmd, err := mapper.ToSetCustomerBalanceCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.UpdateBalance(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) ModifyBalance(ctx context.Context, request *customerv1.ModifyCustomerBalanceRequest) (*customerv1.ModifyCustomerBalanceResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.ModifyBalance")
  defer span.End()

  cmd, err := mapper.ToModifyCustomerBalanceCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.UpdateBalance(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) UpdatePassword(ctx context.Context, request *customerv1.UpdateCustomerPasswordRequest) (*customerv1.UpdateCustomerPasswordResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.UpdatePassword")
  defer span.End()

  cmd, err := mapper.ToUpdateCustomerPasswordCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.UpdatePassword(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) UpdatePhoto(ctx context.Context, request *customerv1.UpdateCustomerPhotoRequest) (*customerv1.UpdateCustomerPhotoResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.UpdatePhoto")
  defer span.End()

  cmd, err := mapper.ToUpdateCustomerPhotoCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.UpdatePhoto(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) UpdateVoucher(ctx context.Context, request *customerv1.UpdateCustomerVoucherRequest) (*customerv1.UpdateCustomerVoucherResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.UpdateVoucher")
  defer span.End()

  cmd, err := mapper.ToUpdateCustomerVoucherCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.UpdateVoucher(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) VerifyEmail(ctx context.Context, request *customerv1.VerifyCustomerEmailRequest) (*customerv1.VerifyCustomerEmailResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.VerifyEmail")
  defer span.End()

  cmd, err := mapper.ToVerifyCustomerEmailCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.VerifyEmail(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) VerifyEmailInstantiate(ctx context.Context, request *customerv1.VerifyCustomerEmailInstantiateRequest) (*customerv1.VerifyCustomerEmailInstantiateResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.VerifyEmailInstantiate")
  defer span.End()

  cmd, err := mapper.ToVerifyCustomerEmailRequestCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.VerifyEmailRequest(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) Create(ctx context.Context, request *customerv1.CreateCustomerRequest) (*customerv1.CreateCustomerResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.Create")
  defer span.End()

  command, err := mapper.ToCreateCustomerCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, stat := c.svc.Create(ctx, &command)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &customerv1.CreateCustomerResponse{
    CustomerId: result.String(),
  }
  return resp, nil
}

func (c *CustomerCommandHandler) Update(ctx context.Context, request *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.Update")
  defer span.End()

  command, err := mapper.ToUpdateCustomerCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.Update(ctx, &command)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) Disable(ctx context.Context, request *customerv1.DisableCustomerRequest) (*customerv1.DisableCustomerResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.Disable")
  defer span.End()

  command, err := mapper.ToDisableCustomerCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.Disable(ctx, &command)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CustomerCommandHandler) Enable(ctx context.Context, request *customerv1.EnableCustomerRequest) (*customerv1.EnableCustomerResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerCommandHandler.Enable")
  defer span.End()

  command, err := mapper.ToEnableCustomerCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.Enable(ctx, &command)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
