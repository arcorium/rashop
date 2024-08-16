package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/internal/domain/event"
  "rashop/services/customer/internal/domain/repository"
)

type ICustomerEmailVerifiedConsumer interface {
  handler.Consumer[*event.CustomerEmailVerifiedV1]
}

func NewCustomerEmailVerifiedConsumer(repo repository.ICustomer) ICustomerEmailVerifiedConsumer {
  return &customerEmailVerifiedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerEmailVerifiedConsumer struct {
  commonHandlerField
}

func (c *customerEmailVerifiedConsumer) Handle(ctx context.Context, e *event.CustomerEmailVerifiedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerEmailVerifiedConsumer.Handle")
  defer span.End()

  stat := HandleSimple(ctx, e.CustomerId, c.repo, e)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
  }
  return stat
}
