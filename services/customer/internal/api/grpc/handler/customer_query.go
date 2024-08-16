package handler

import (
  "context"
  "github.com/arcorium/rashop/proto/gen/go/customer/v1"
  "github.com/arcorium/rashop/shared/logger"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "rashop/services/customer/internal/api/grpc/mapper"
  "rashop/services/customer/internal/app/query"
  "rashop/services/customer/internal/app/service"
  "rashop/services/customer/pkg/tracer"
)

func NewCustomerQuery(svc service.ICustomerQuery) CustomerQueryHandler {
  return CustomerQueryHandler{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type CustomerQueryHandler struct {
  customerv1.UnimplementedCustomerQueryServiceServer

  svc    service.ICustomerQuery
  tracer trace.Tracer
}

func (c *CustomerQueryHandler) Register(server *grpc.Server) {
  customerv1.RegisterCustomerQueryServiceServer(server, c)
}

func (c *CustomerQueryHandler) Get(ctx context.Context, request *customerv1.GetCustomersRequest) (*customerv1.GetCustomerResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryHandler.GetCustomers")
  defer span.End()

  dtos := mapper.ToPagedParameter(request)
  logger.Infof("Value: %v", dtos)

  customers, stat := c.svc.GetCustomers(ctx, &query.GetCustomersQuery{PagedElementDTO: dtos})
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &customerv1.GetCustomerResponse{
    Customers: sharedUtil.CastSliceP(customers.Data, mapper.ToProtoCustomer),
    Detail:    mapper.ToProtoPagedElementDetails(&customers),
  }
  return resp, nil
}

func (c *CustomerQueryHandler) FindByIds(ctx context.Context, request *customerv1.FindCustomerByIdsRequest) (*customerv1.FindCustomerByIdsResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryHandler.FindByIds")
  defer span.End()

  query, err := mapper.ToGetCustomerByIdsQuery(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, stat := c.svc.GetCustomerByIds(ctx, &query)
  if stat.IsError() {
    spanUtil.RecordError(err, span)
    return nil, stat.ToGRPCError()
  }

  resp := &customerv1.FindCustomerByIdsResponse{
    Customers: sharedUtil.CastSliceP(result, mapper.ToProtoCustomer),
  }
  return resp, nil
}

func (c *CustomerQueryHandler) FindAddresses(ctx context.Context, request *customerv1.FindCustomerAddressesRequest) (*customerv1.FindCustomerAddressesResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryHandler.FindAddresses")
  defer span.End()

  query, err := mapper.ToGetCustomerAddressesQuery(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, stat := c.svc.GetCustomerAddresses(ctx, &query)
  if stat.IsError() {
    spanUtil.RecordError(err, span)
    return nil, stat.ToGRPCError()
  }

  resp := &customerv1.FindCustomerAddressesResponse{
    Addresses: sharedUtil.CastSliceP(result, mapper.ToProtoAddress),
  }
  return resp, nil
}

func (c *CustomerQueryHandler) FindVouchers(ctx context.Context, request *customerv1.FindCustomerVouchersRequest) (*customerv1.FindCustomerVouchersResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CustomerQueryHandler.FindVouchers")
  defer span.End()

  query, err := mapper.ToGetCustomerVouchersQuery(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, stat := c.svc.GetCustomerVouchers(ctx, &query)
  if stat.IsError() {
    spanUtil.RecordError(err, span)
    return nil, stat.ToGRPCError()
  }

  resp := &customerv1.FindCustomerVouchersResponse{
    Vouchers: sharedUtil.CastSliceP(result, mapper.ToProtoVoucher),
  }
  return resp, nil
}
