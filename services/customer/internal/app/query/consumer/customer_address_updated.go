package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/internal/domain/event"
  "rashop/services/customer/internal/domain/repository"
)

type ICustomerAddressUpdatedConsumer interface {
  handler.Consumer[*event.CustomerAddressUpdatedV1]
}

func NewCustomerAddressUpdatedConsumer(repo repository.ICustomer) ICustomerAddressUpdatedConsumer {
  return &customerAddressUpdatedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerAddressUpdatedConsumer struct {
  commonHandlerField
}

func (c *customerAddressUpdatedConsumer) Handle(ctx context.Context, e *event.CustomerAddressUpdatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerAddressUpdatedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
