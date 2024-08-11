package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/internal/domain/repository"
)

type ICustomerVoucherUpdatedConsumer interface {
  handler.Consumer[*event.CustomerVoucherUpdatedV1]
}

func NewCustomerVoucherUpdatedConsumer(repo repository.ICustomer) ICustomerVoucherUpdatedConsumer {
  return &customerVoucherUpdatedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerVoucherUpdatedConsumer struct {
  commonHandlerField
}

func (c *customerVoucherUpdatedConsumer) Handle(ctx context.Context, e *event.CustomerVoucherUpdatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerVoucherUpdatedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
