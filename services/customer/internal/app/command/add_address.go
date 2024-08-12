package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
)

type IAddCustomerAddressHandler interface {
  handler.Command[*AddCustomerAddressCommand, types.Id]
}

func NewAddCustomerAddressHandler(parameter cqrs.CommonHandlerParameter) IAddCustomerAddressHandler {
  return &addCustomerAddressHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type addCustomerAddressHandler struct {
  basicHandler
}

func (a *addCustomerAddressHandler) Handle(ctx context.Context, cmd *AddCustomerAddressCommand) (types.Id, status.Object) {
  ctx, span := a.tracer.Start(ctx, "addCustomerAddressHandler.Handle")
  defer span.End()

  address, err := cmd.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  // Get aggregate
  customers, err := a.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository(err)
  }
  customer := &customers[0]

  // Update
  ev := customer.AddAddress(&address)
  err = customer.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrBadRequest(err)
  }

  // Update aggregate
  err = a.repo.Update(ctx, customer)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository(err)
  }

  // Forward all domain events
  err = a.publisher.PublishAggregate(ctx, customer)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  return address.Id, status.Created()
}
