package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
)

type IUpdateCustomerAddressHandler interface {
  handler.CommandUnit[*UpdateCustomerAddressCommand]
}

func NewUpdateCustomerAddressHandler(parameter cqrs.CommonHandlerParameter) IUpdateCustomerAddressHandler {
  return &updateCustomerAddressHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type updateCustomerAddressHandler struct {
  basicHandler
}

func (a *updateCustomerAddressHandler) Handle(ctx context.Context, cmd *UpdateCustomerAddressCommand) status.Object {
  ctx, span := a.tracer.Start(ctx, "updateCustomerAddressHandler.Handle")
  defer span.End()

  // Get aggregate
  customers, err := a.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }
  current := &customers[0]

  // Update address
  ev, err := current.UpdateAddress(cmd.AddressId, cmd.ToPredicate())
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

  // Forward all domain events
  err = a.publisher.PublishAggregate(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Success()
}
