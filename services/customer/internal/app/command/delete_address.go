package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
)

type IDeleteCustomerAddressHandler interface {
  handler.CommandUnit[*DeleteCustomerAddressCommand]
}

func NewDeleteCustomerAddressHandler(parameter cqrs.CommonHandlerParameter) IDeleteCustomerAddressHandler {
  return &deleteCustomerAddressHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type deleteCustomerAddressHandler struct {
  basicHandler
}

func (a *deleteCustomerAddressHandler) Handle(ctx context.Context, cmd *DeleteCustomerAddressCommand) status.Object {
  ctx, span := a.tracer.Start(ctx, "deleteCustomerAddressHandler.Handle")
  defer span.End()

  // Get aggregate
  customers, err := a.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }
  customer := &customers[0]

  // Delete address
  ev, err := customer.DeleteAddress(cmd.AddressId)
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
