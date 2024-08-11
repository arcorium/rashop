package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
)

type ISetCustomerDefaultAddressHandler interface {
  handler.CommandUnit[*SetCustomerDefaultAddressCommand]
}

func NewSetCustomerDefaultAddressHandler(parameter cqrs.CommonHandlerParameter) ISetCustomerDefaultAddressHandler {
  return &setCustomerDefaultAddressHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type setCustomerDefaultAddressHandler struct {
  basicHandler
}

func (u *setCustomerDefaultAddressHandler) Handle(ctx context.Context, cmd *SetCustomerDefaultAddressCommand) status.Object {
  ctx, span := u.tracer.Start(ctx, "setCustomerDefaultAddressHandler.Handle")
  defer span.End()

  customers, err := u.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &customers[0]
  ev, err := current.SetDefaultAddress(cmd.AddressId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err = current.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err = u.repo.Update(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  err = u.publisher.PublishAggregate(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Updated()
}
