package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
)

type IDeleteCustomerVoucherHandler interface {
  handler.CommandUnit[*DeleteCustomerVoucherCommand]
}

func NewDeleteCustomerVoucherHandler(parameter cqrs.CommonHandlerParameter) IDeleteCustomerVoucherHandler {
  return &deleteCustomerVoucherHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type deleteCustomerVoucherHandler struct {
  basicHandler
}

func (a *deleteCustomerVoucherHandler) Handle(ctx context.Context, cmd *DeleteCustomerVoucherCommand) status.Object {
  ctx, span := a.tracer.Start(ctx, "deleteCustomerVoucherHandler.Handle")
  defer span.End()

  // Get aggregate
  customers, err := a.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }
  customer := &customers[0]

  // Delete voucher
  ev, err := customer.DeleteVoucher(cmd.VoucherId)
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
  err = a.publisher.PublishAggregate(ctx, customer)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Success()
}
