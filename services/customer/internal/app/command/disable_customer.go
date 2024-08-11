package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
)

type IDisableCustomerHandler interface {
  handler.CommandUnit[*DisableCustomerCommand]
}

func NewDisableCustomerHandler(parameter cqrs.CommonHandlerParameter) IDisableCustomerHandler {
  return &disableCustomerHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type disableCustomerHandler struct {
  basicHandler
}

func (d *disableCustomerHandler) Handle(ctx context.Context, command *DisableCustomerCommand) status.Object {
  ctx, span := d.tracer.Start(ctx, "disableCustomerHandler.Handle")
  defer span.End()

  customers, err := d.repo.FindByIds(ctx, command.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &customers[0]
  ev, err := current.Disable()
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
  err = d.repo.Update(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  err = d.publisher.PublishAggregate(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Updated()
}
