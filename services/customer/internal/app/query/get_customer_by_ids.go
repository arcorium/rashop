package query

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "go.opentelemetry.io/otel/trace"
  "rashop/services/customer/internal/app/dto"
  "rashop/services/customer/internal/app/mapper"
  "rashop/services/customer/internal/domain/repository"
  "rashop/services/customer/pkg/tracer"
)

type IGetCustomerByIdsHandler interface {
  handler.Command[*GetCustomerByIdsQuery, []dto.CustomerResponseDTO]
}

func NewGetCustomerByIdHandler(customer repository.ICustomer) IGetCustomerByIdsHandler {
  return &getCustomerByIdsHandler{
    repo:   customer,
    tracer: tracer.Get(),
  }
}

type getCustomerByIdsHandler struct {
  repo   repository.ICustomer
  tracer trace.Tracer
}

func (h *getCustomerByIdsHandler) Handle(ctx context.Context, query *GetCustomerByIdsQuery) ([]dto.CustomerResponseDTO, status.Object) {
  ctx, span := h.tracer.Start(ctx, "getCustomerAddresses.Handle")
  defer span.End()

  customers, err := h.repo.FindByIds(ctx, query.CustomerIds...)
  if err != nil {
    span.RecordError(err)
    return nil, status.FromRepository(err)
  }

  resp := sharedUtil.CastSliceP(customers, mapper.ToCustomerResponseDTO)
  return resp, status.Succeed()
}
