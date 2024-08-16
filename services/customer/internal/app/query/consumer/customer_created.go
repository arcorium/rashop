package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/internal/domain/entity"
  "rashop/services/customer/internal/domain/event"
  "rashop/services/customer/internal/domain/repository"
)

type ICustomerCreatedConsumer interface {
  handler.Consumer[*event.CustomerCreatedV1]
}

func NewCustomerCreatedConsumer(repo repository.ICustomer) ICustomerCreatedConsumer {
  return &customerCreatedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerCreatedConsumer struct {
  commonHandlerField
}

func (c *customerCreatedConsumer) Handle(ctx context.Context, e *event.CustomerCreatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerCreatedConsumer.Handle")
  defer span.End()

  customer := entity.Customer{}
  err := customer.ApplyEvent(e)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  err = c.repo.Create(ctx, &customer)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  return status.Created()
}
