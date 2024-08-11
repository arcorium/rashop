package command

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
)

type IUpdateCustomerHandler interface {
  handler.CommandUnit[*UpdateCustomerCommand]
}

func NewUpdateCustomerHandler(parameter cqrs.CommonHandlerParameter) IUpdateCustomerHandler {
  return &updateCustomerHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type updateCustomerHandler struct {
  basicHandler
}

func (u *updateCustomerHandler) Handle(ctx context.Context, cmd *UpdateCustomerCommand) status.Object {
  ctx, span := u.tracer.Start(ctx, "updateCustomerHandler.Handle")
  defer span.End()

  customers, err := u.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &customers[0]
  ev := current.Update(cmd.ToPredicate())
  err = current.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  // Publish email changed integration event
  if current.IsEmailChanged() {
    current.AddEvents(intev.NewCustomerEmailChangedV1(current.Id, current.Email))
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
