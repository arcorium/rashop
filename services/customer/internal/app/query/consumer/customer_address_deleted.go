package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/internal/domain/repository"
)

type ICustomerAddressDeletedConsumer interface {
  handler.Consumer[*event.CustomerAddressDeletedV1]
}

func NewCustomerAddressDeletedConsumer(repo repository.ICustomer) ICustomerAddressDeletedConsumer {
  return &customerAddressDeletedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerAddressDeletedConsumer struct {
  commonHandlerField
}

func (c *customerAddressDeletedConsumer) Handle(ctx context.Context, e *event.CustomerAddressDeletedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerAddressDeletedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
