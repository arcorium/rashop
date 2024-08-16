package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/internal/domain/event"
  "rashop/services/customer/internal/domain/repository"
)

type ICustomerDefaultAddressSetConsumer interface {
  handler.Consumer[*event.CustomerDefaultAddressUpdatedV1]
}

func NewCustomerDefaultAddressSetConsumer(repo repository.ICustomer) ICustomerDefaultAddressSetConsumer {
  return &customerDefaultAddressSetConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerDefaultAddressSetConsumer struct {
  commonHandlerField
}

func (c *customerDefaultAddressSetConsumer) Handle(ctx context.Context, e *event.CustomerDefaultAddressUpdatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerDefaultAddressSetConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
