package query

import (
  "context"
  sharedDto "github.com/arcorium/rashop/shared/dto"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "go.opentelemetry.io/otel/trace"
  "rashop/services/customer/internal/app/dto"
  "rashop/services/customer/internal/app/mapper"
  "rashop/services/customer/internal/domain/repository"
  "rashop/services/customer/pkg/tracer"
)

type IGetCustomersHandler interface {
  handler.Command[*GetCustomersQuery, sharedDto.PagedElementResult[dto.CustomerResponseDTO]]
}

func NewGetCustomersHandler(customer repository.ICustomer) IGetCustomersHandler {
  return &getCustomersHandler{
    repo:   customer,
    tracer: tracer.Get(),
  }
}

type getCustomersHandler struct {
  repo   repository.ICustomer
  tracer trace.Tracer
}

func (h *getCustomersHandler) Handle(ctx context.Context, query *GetCustomersQuery) (sharedDto.PagedElementResult[dto.CustomerResponseDTO], status.Object) {
  ctx, span := h.tracer.Start(ctx, "getCustomerAddresses.Handle")
  defer span.End()

  result, err := h.repo.Get(ctx, query.ToQueryParam())
  if err != nil {
    span.RecordError(err)
    return sharedDto.PagedElementResult[dto.CustomerResponseDTO]{}, status.FromRepository(err)
  }

  customers := sharedUtil.CastSliceP(result.Data, mapper.ToCustomerResponseDTO)
  resp := sharedDto.NewPagedElementResult2(customers, &query.PagedElementDTO, result.Total)
  return resp, status.Succeed()
}
