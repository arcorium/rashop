package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/internal/domain/event"
  "rashop/services/customer/internal/domain/repository"
)

type ICustomerPhotoUpdatedConsumer interface {
  handler.Consumer[*event.CustomerPhotoUpdatedV1]
}

func NewCustomerPhotoUpdatedConsumer(repo repository.ICustomer) ICustomerPhotoUpdatedConsumer {
  return &customerPhotoUpdatedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerPhotoUpdatedConsumer struct {
  commonHandlerField
}

func (c *customerPhotoUpdatedConsumer) Handle(ctx context.Context, e *event.CustomerPhotoUpdatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerPhotoUpdatedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
