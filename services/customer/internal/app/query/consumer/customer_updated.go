package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/internal/domain/repository"
)

type ICustomerUpdatedConsumer interface {
  handler.Consumer[*event.CustomerUpdatedV1]
}

func NewCustomerUpdatedConsumer(repo repository.ICustomer) ICustomerUpdatedConsumer {
  return &customerUpdatedConsumer{
    commonHandlerField: newBasicHandler(repo),
  }
}

type customerUpdatedConsumer struct {
  commonHandlerField
}

func (c *customerUpdatedConsumer) Handle(ctx context.Context, e *event.CustomerUpdatedV1) status.Object {
  ctx, span := c.tracer.Start(ctx, "customerUpdatedConsumer.Handle")
  defer span.End()

  custId, err := types.IdFromString(e.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  customers, err := c.repo.FindByIds(ctx, custId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &customers[0]
  err = current.ApplyEvent(e)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  err = c.repo.Update(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  return status.Success()
}
