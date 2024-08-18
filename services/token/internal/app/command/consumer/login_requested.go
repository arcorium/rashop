package consumer

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type ILoginRequestedHandler interface {
  handler.Consumer[*intev.LoginTokenRequestedV1]
}

func NewLoginRequestedHandler(parameter CommonHandlerParameter) ILoginRequestedHandler {
  return &loginRequestedHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type loginRequestedHandler struct {
  commonHandler
}

func (e *loginRequestedHandler) Handle(ctx context.Context, ev *intev.LoginTokenRequestedV1) status.Object {
  ctx, span := e.tracer.Start(ctx, "loginRequestedHandler.Handle")
  defer span.End()

  // Map into command
  cmd, err := MapLoginTokenRequested(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Generate token
  _, stat := e.generate.Handle(ctx, &cmd)
  return stat
}
