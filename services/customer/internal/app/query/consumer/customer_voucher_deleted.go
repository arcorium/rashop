package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/internal/domain/repository"
)

type ICustomerVoucherDeletedConsumer interface {
  handler.Consumer[*event.CustomerVoucherDeletedV1]
}

func NewCustomerVoucherDeletedConsumer(repo repository.ICustomer) ICustomerVoucherDeletedConsumer {
  return &customerVoucherDeletedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerVoucherDeletedConsumer struct {
  commonHandlerField
}

func (c *customerVoucherDeletedConsumer) Handle(ctx context.Context, e *event.CustomerVoucherDeletedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerVoucherDeletedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
