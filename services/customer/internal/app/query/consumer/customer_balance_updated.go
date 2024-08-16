package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/internal/domain/event"
  "rashop/services/customer/internal/domain/repository"
)

type ICustomerBalanceUpdatedConsumer interface {
  handler.Consumer[*event.CustomerBalanceUpdatedV1]
}

func NewCustomerBalanceUpdatedConsumer(repo repository.ICustomer) ICustomerBalanceUpdatedConsumer {
  return &customerBalanceUpdatedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerBalanceUpdatedConsumer struct {
  commonHandlerField
}

func (c *customerBalanceUpdatedConsumer) Handle(ctx context.Context, e *event.CustomerBalanceUpdatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerBalanceUpdatedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
