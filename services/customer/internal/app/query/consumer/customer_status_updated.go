package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/internal/domain/repository"
)

type ICustomerStatusUpdatedConsumer interface {
  handler.Consumer[*event.CustomerStatusUpdatedV1]
}

func NewCustomerStatusUpdatedConsumer(repo repository.ICustomer) ICustomerStatusUpdatedConsumer {
  return &customerStatusUpdatedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerStatusUpdatedConsumer struct {
  commonHandlerField
}

func (c *customerStatusUpdatedConsumer) Handle(ctx context.Context, e *event.CustomerStatusUpdatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerStatusUpdatedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
