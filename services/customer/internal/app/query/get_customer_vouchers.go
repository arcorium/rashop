package query

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "rashop/services/customer/internal/app/dto"
  "rashop/services/customer/internal/app/mapper"
  "rashop/services/customer/internal/domain/repository"
  "rashop/services/customer/pkg/tracer"
)

type IGetCustomerVouchersHandler interface {
  handler.Command[*GetCustomerVouchersQuery, []dto.VoucherResponseDTO]
}

func NewGetCustomerVouchers(customer repository.ICustomer) IGetCustomerVouchersHandler {
  return &getCustomerVouchers{
    repo:   customer,
    tracer: tracer.Get(),
  }
}

type getCustomerVouchers struct {
  repo   repository.ICustomer
  tracer trace.Tracer
}

func (g *getCustomerVouchers) Handle(ctx context.Context, query *GetCustomerVouchersQuery) ([]dto.VoucherResponseDTO, status.Object) {
  ctx, span := g.tracer.Start(ctx, "getCustomerVouchers.Handle")
  defer span.End()

  customers, err := g.repo.FindByIds(ctx, query.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err)
  }

  current := &customers[0]
  result := sharedUtil.CastSliceP(current.Vouchers.Elements(), mapper.ToVoucherResponseDTO)
  return result, status.Succeed()
}
