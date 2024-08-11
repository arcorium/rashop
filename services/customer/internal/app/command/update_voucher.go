package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
)

type IUpdateCustomerVoucherHandler interface {
  handler.CommandUnit[*UpdateCustomerVoucherCommand]
}

func NewUpdateCustomerVoucherHandler(parameter cqrs.CommonHandlerParameter) IUpdateCustomerVoucherHandler {
  return &updateCustomerVoucherHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type updateCustomerVoucherHandler struct {
  basicHandler
}

func (a *updateCustomerVoucherHandler) Handle(ctx context.Context, cmd *UpdateCustomerVoucherCommand) status.Object {
  ctx, span := a.tracer.Start(ctx, "updateCustomerVoucherHandler.Handle")
  defer span.End()

  // Get aggregate
  customers, err := a.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }
  current := &customers[0]

  // Update voucher
  ev, err := current.UpdateVoucher(cmd.VoucherId, cmd.ToPredicate())
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err = current.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  // Update aggregate
  err = a.repo.Update(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  // Publish all domain events
  err = a.publisher.PublishAggregate(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Success()
}
