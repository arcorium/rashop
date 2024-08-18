package consumer

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IResetPasswordRequestedHandler interface {
  handler.Consumer[*intev.CustomerResetPasswordRequestedV1]
}

func NewResetPasswordRequestedHandler(parameter CommonHandlerParameter) IResetPasswordRequestedHandler {
  return &resetPasswordRequestedHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type resetPasswordRequestedHandler struct {
  commonHandler
}

func (e *resetPasswordRequestedHandler) Handle(ctx context.Context, ev *intev.CustomerResetPasswordRequestedV1) status.Object {
  ctx, span := e.tracer.Start(ctx, "resetPasswordRequestedHandler.Handle")
  defer span.End()

  // Map into command
  cmd, err := MapResetPasswordRequested(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Generate token
  _, stat := e.generate.Handle(ctx, &cmd)
  return stat
}
