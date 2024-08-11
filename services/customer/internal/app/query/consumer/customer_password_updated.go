package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/internal/domain/repository"
)

type ICustomerPasswordUpdatedConsumer interface {
  handler.Consumer[*event.CustomerPasswordUpdatedV1]
}

func NewCustomerPasswordUpdatedConsumer(repo repository.ICustomer) ICustomerPasswordUpdatedConsumer {
  return &customerPasswordUpdatedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerPasswordUpdatedConsumer struct {
  commonHandlerField
}

func (c *customerPasswordUpdatedConsumer) Handle(ctx context.Context, e *event.CustomerPasswordUpdatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerPasswordUpdatedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
