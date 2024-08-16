package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/pkg/cqrs"
)

type IAddCustomerVouchersHandler interface {
  handler.CommandUnit[*AddCustomerVoucherCommand]
}

func NewAddCustomerVouchersHandler(parameter cqrs.CommonHandlerParameter) IAddCustomerVouchersHandler {
  return &addCustomerVoucherHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type addCustomerVoucherHandler struct {
  basicHandler
}

func (a *addCustomerVoucherHandler) Handle(ctx context.Context, cmd *AddCustomerVoucherCommand) status.Object {
  ctx, span := a.tracer.Start(ctx, "deleteCustomerVoucherHandler.Handle")
  defer span.End()

  voucher := cmd.ToDomain()
  // GetCustomers aggregate
  customers, err := a.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }
  customer := &customers[0]

  // Add voucher
  ev, err := customer.AddVoucher(&voucher)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err = customer.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  // Update aggregate
  err = a.repo.Update(ctx, customer)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  // Forward all domain events
  err = a.publisher.Publish(ctx, customer)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Succeed()
}
