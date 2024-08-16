package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/pkg/cqrs"
)

type IResetPasswordHandler interface {
  handler.CommandUnit[*ResetCustomerPasswordCommand]
}

func NewResetPasswordHandler(parameter cqrs.CommonHandlerParameter) IResetPasswordHandler {
  return &resetPasswordHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type resetPasswordHandler struct {
  basicHandler
}

func (r *resetPasswordHandler) Handle(ctx context.Context, cmd *ResetCustomerPasswordCommand) status.Object {
  ctx, span := r.tracer.Start(ctx, "resetPasswordHandler.Handle")
  defer span.End()

  // TODO: Validate token
  custId := types.NullId()

  customers, err := r.repo.FindByIds(ctx, custId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &customers[0]
  ev, err := current.ResetPassword(cmd.NewPassword)
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
  err = r.repo.Update(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  err = r.publisher.Publish(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Updated()
}
