package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/internal/domain/event"
  "rashop/services/customer/internal/domain/repository"
)

type ICustomerAddressAddedConsumer interface {
  handler.Consumer[*event.CustomerAddressAddedV1]
}

func NewCustomerAddressAddedConsumer(repo repository.ICustomer) ICustomerAddressAddedConsumer {
  return &customerAddressAddedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerAddressAddedConsumer struct {
  commonHandlerField
}

func (c *customerAddressAddedConsumer) Handle(ctx context.Context, e *event.CustomerAddressAddedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerAddressAddedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
