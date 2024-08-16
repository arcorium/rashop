package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/internal/domain/event"
  "rashop/services/customer/internal/domain/repository"
)

type ICustomerVoucherAddedConsumer interface {
  handler.Consumer[*event.CustomerVoucherAddedV1]
}

func NewCustomerVoucherAddedConsumer(repo repository.ICustomer) ICustomerVoucherAddedConsumer {
  return &customerVouchersAddedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerVouchersAddedConsumer struct {
  commonHandlerField
}

func (c *customerVouchersAddedConsumer) Handle(ctx context.Context, e *event.CustomerVoucherAddedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerVouchersAddedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
