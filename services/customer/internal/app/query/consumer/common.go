package consumer

import (
  "context"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  "go.opentelemetry.io/otel/trace"
  "rashop/services/customer/internal/domain/repository"
  "rashop/services/customer/pkg/tracer"
)

func newBasicHandler(repo repository.ICustomer, opts ...trace.TracerOption) commonHandlerField {
  return commonHandlerField{
    repo:   repo,
    tracer: tracer.Get(opts...),
  }
}

type commonHandlerField struct {
  repo   repository.ICustomer
  tracer trace.Tracer
}

// HandleSimple is common handler for consumer handler. It will get the aggregate based on the 'customerId',
// apply 'ev' event and update the aggregate into repository
func HandleSimple[E types.Event](ctx context.Context, customerId string, repo repository.ICustomer, ev E) status.Object {
  // Check id
  custId, err := types.IdFromString(customerId)
  if err != nil {
    return status.ErrInternal(err)
  }

  // GetCustomers the aggregate based on the id
  customers, err := repo.FindByIds(ctx, custId)
  if err != nil {
    return status.FromRepository(err)
  }

  // Apply the event
  current := &customers[0]
  err = current.ApplyEvent(ev)
  if err != nil {
    return status.ErrBadRequest(err)
  }

  // Update modified aggregate
  err = repo.Update(ctx, current)
  if err != nil {
    return status.FromRepository(err)
  }

  // All consumer will only return status.SUCCESS
  return status.Succeed()
}
