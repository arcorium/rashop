package consumer

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type ICustomerCreatedHandler interface {
  handler.Consumer[*intev.CustomerCreatedV1]
}

func NewCustomerCreatedHandler(parameter CommonHandlerParameter) ICustomerCreatedHandler {
  return &customerCreatedHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type customerCreatedHandler struct {
  commonHandler
}

func (e *customerCreatedHandler) Handle(ctx context.Context, ev *intev.CustomerCreatedV1) status.Object {
  ctx, span := e.tracer.Start(ctx, "customerCreatedHandler.Handle")
  defer span.End()

  // Map into command
  cmd, err := MapCustomerCreated(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Generate token
  _, stat := e.generate.Handle(ctx, &cmd)
  return stat
}
