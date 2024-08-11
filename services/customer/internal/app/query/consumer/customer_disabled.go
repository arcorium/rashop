package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/internal/domain/repository"
)

type ICustomerDisabledConsumer interface {
  handler.Consumer[*event.CustomerStatusUpdatedV1]
}

func NewCustomerDisabledConsumer(repo repository.ICustomer) ICustomerDisabledConsumer {
  return &customerDisabledConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerDisabledConsumer struct {
  commonHandlerField
}

func (c *customerDisabledConsumer) Handle(ctx context.Context, e *event.CustomerStatusUpdatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerDisabledConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
