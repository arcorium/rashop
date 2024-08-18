package consumer

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IEmailVerificationRequestedHandler interface {
  handler.Consumer[*intev.CustomerEmailVerificationRequestedV1]
}

func NewEmailVerificationRequestedHandler(parameter CommonHandlerParameter) IEmailVerificationRequestedHandler {
  return &emailVerificationRequestedHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type emailVerificationRequestedHandler struct {
  commonHandler
}

func (e *emailVerificationRequestedHandler) Handle(ctx context.Context, ev *intev.CustomerEmailVerificationRequestedV1) status.Object {
  ctx, span := e.tracer.Start(ctx, "emailVerificationRequestedHandler.Handle")
  defer span.End()

  // Map into command
  cmd, err := MapCustomerEmailVerificationRequested(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Generate token
  _, stat := e.generate.Handle(ctx, &cmd)
  return stat
}
