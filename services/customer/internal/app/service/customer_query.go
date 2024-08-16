package service

import (
  "context"
  sharedDto "github.com/arcorium/rashop/shared/dto"
  "github.com/arcorium/rashop/shared/status"
  "rashop/services/customer/internal/app/dto"
  "rashop/services/customer/internal/app/query"
  "rashop/services/customer/internal/domain/repository"
)

type ICustomerQuery interface {
  GetCustomers(ctx context.Context, input *query.GetCustomersQuery) (sharedDto.PagedElementResult[dto.CustomerResponseDTO], status.Object)
  GetCustomerByIds(ctx context.Context, input *query.GetCustomerByIdsQuery) ([]dto.CustomerResponseDTO, status.Object)
  GetCustomerAddresses(ctx context.Context, input *query.GetCustomerAddressesQuery) ([]dto.AddressResponseDTO, status.Object)
  GetCustomerVouchers(ctx context.Context, input *query.GetCustomerVouchersQuery) ([]dto.VoucherResponseDTO, status.Object)
}

func NewCustomerQuery(config CustomerQueryFactory) ICustomerQuery {
  return &customerQueryService{
    CustomerQueryFactory: config,
  }
}

func DefaultCustomerQueryFactory(repo repository.ICustomer) CustomerQueryFactory {
  return CustomerQueryFactory{
    Get:          query.NewGetCustomersHandler(repo),
    GetById:      query.NewGetCustomerByIdHandler(repo),
    GetAddresses: query.NewGetCustomerAddresses(repo),
    GetVouchers:  query.NewGetCustomerVouchers(repo),
  }
}

type CustomerQueryFactory struct {
  Get          query.IGetCustomersHandler
  GetById      query.IGetCustomerByIdsHandler
  GetAddresses query.IGetCustomerAddressesHandler
  GetVouchers  query.IGetCustomerVouchersHandler
}

type customerQueryService struct {
  CustomerQueryFactory
}

func (c *customerQueryService) GetCustomers(ctx context.Context, input *query.GetCustomersQuery) (sharedDto.PagedElementResult[dto.CustomerResponseDTO], status.Object) {
  return c.Get.Handle(ctx, input)
}

func (c *customerQueryService) GetCustomerByIds(ctx context.Context, input *query.GetCustomerByIdsQuery) ([]dto.CustomerResponseDTO, status.Object) {
  return c.GetById.Handle(ctx, input)
}

func (c *customerQueryService) GetCustomerAddresses(ctx context.Context, input *query.GetCustomerAddressesQuery) ([]dto.AddressResponseDTO, status.Object) {
  return c.GetAddresses.Handle(ctx, input)
}

func (c *customerQueryService) GetCustomerVouchers(ctx context.Context, input *query.GetCustomerVouchersQuery) ([]dto.VoucherResponseDTO, status.Object) {
  return c.GetVouchers.Handle(ctx, input)
}
