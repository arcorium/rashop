package command

import (
  "context"
  "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/internal/domain/entity"
  "mini-shop/services/user/pkg/cqrs"
)

type ICreateCustomerHandler interface {
  handler.Command[*CreateCustomerCommand, types.Id]
}

func NewCreateCustomerHandler(parameter cqrs.CommonHandlerParameter) ICreateCustomerHandler {
  return &createCustomerHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type createCustomerHandler struct {
  basicHandler
}

func (c *createCustomerHandler) Handle(ctx context.Context, command *CreateCustomerCommand) (types.Id, status.Object) {
  ctx, span := c.tracer.Start(ctx, "createCustomerHandler.Handle")
  defer span.End()

  customer, err := command.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  customer, ev := entity.CreateCustomer(&customer)
  err = customer.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrBadRequest(err)
  }

  // Create persistent
  err = c.repo.Create(ctx, &customer)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepositoryExist(err)
  }

  // Add integration event
  integrationEv := intev.NewCustomerCreatedV1(customer.Id, customer.Email, customer.Name.User)
  customer.AddEvents(integrationEv)

  // Forward all domain events
  err = c.publisher.PublishAggregate(ctx, &customer)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  return customer.Id, status.Created()
}
