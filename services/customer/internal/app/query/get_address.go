package query

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "go.opentelemetry.io/otel/trace"
  "mini-shop/services/user/internal/app/dto"
  "mini-shop/services/user/internal/app/mapper"
  "mini-shop/services/user/internal/domain/repository"
  "mini-shop/services/user/pkg/tracer"
)

type IGetCustomerAddressesHandler interface {
  handler.Command[*GetCustomerAddressesQuery, []dto.AddressResponseDTO]
}

func NewGetCustomerAddresses(customer repository.ICustomer) IGetCustomerAddressesHandler {
  return &getCustomerAddresses{
    repo:   customer,
    tracer: tracer.Get(),
  }
}

type getCustomerAddresses struct {
  repo   repository.ICustomer
  tracer trace.Tracer
}

func (h *getCustomerAddresses) Handle(ctx context.Context, query *GetCustomerAddressesQuery) ([]dto.AddressResponseDTO, status.Object) {
  ctx, span := h.tracer.Start(ctx, "getCustomerAddresses.Handle")
  defer span.End()

  customers, err := h.repo.FindByIds(ctx, query.CustomerId)
  if err != nil {
    span.RecordError(err)
    return nil, status.FromRepository(err)
  }

  resp := sharedUtil.CastSliceP(customers[0].ShippingAddresses.Elements(), mapper.ToAddressResponseDTO)
  return resp, status.Success()
}
