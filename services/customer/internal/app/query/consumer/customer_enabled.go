package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/internal/domain/repository"
)

type ICustomerEnabledConsumer interface {
  handler.Consumer[*event.CustomerStatusUpdatedV1]
}

func NewCustomerEnabledConsumer(repo repository.ICustomer) ICustomerEnabledConsumer {
  return &customerEnabledConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerEnabledConsumer struct {
  commonHandlerField
}

func (c *customerEnabledConsumer) Handle(ctx context.Context, e *event.CustomerStatusUpdatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerEnabledConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
