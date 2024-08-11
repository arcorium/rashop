package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
)

type IUpdatePasswordHandler interface {
  handler.CommandUnit[*UpdateCustomerPasswordCommand]
}

func NewUpdatePasswordHandler(parameter cqrs.CommonHandlerParameter) IUpdatePasswordHandler {
  return &updatePasswordHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type updatePasswordHandler struct {
  basicHandler
}

func (u *updatePasswordHandler) Handle(ctx context.Context, cmd *UpdateCustomerPasswordCommand) status.Object {
  ctx, span := u.tracer.Start(ctx, "updatePasswordHandler.Handle")
  defer span.End()

  customers, err := u.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &customers[0]
  ev, err := current.UpdatePassword(cmd.LastPassword, cmd.NewPassword)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
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
