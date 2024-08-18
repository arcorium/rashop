package consumer

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type ICustomerEmailChangedHandler interface {
  handler.Consumer[*intev.CustomerEmailUpdatedV1]
}

func NewCustomerEmailChangedHandler(parameter CommonHandlerParameter) ICustomerEmailChangedHandler {
  return &customerEmailChangedHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type customerEmailChangedHandler struct {
  commonHandler
}

func (e *customerEmailChangedHandler) Handle(ctx context.Context, ev *intev.CustomerEmailUpdatedV1) status.Object {
  ctx, span := e.tracer.Start(ctx, "customerEmailChangedHandler.Handle")
  defer span.End()

  // Map into command
  cmd, err := MapCustomerEmailChanged(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Generate token
  _, stat := e.generate.Handle(ctx, &cmd)
  return stat
}
