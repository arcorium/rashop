package service

import (
  "context"
  "github.com/arcorium/rashop/shared/status"
  "mini-shop/services/user/internal/app/dto"
  "mini-shop/services/user/internal/app/query"
  "mini-shop/services/user/internal/domain/repository"
)

type ICustomerQuery interface {
  GetCustomerByIds(ctx context.Context, input *query.GetCustomerByIdsQuery) ([]dto.CustomerResponseDTO, status.Object)
  GetCustomerAddresses(ctx context.Context, input *query.GetCustomerAddressesQuery) ([]dto.AddressResponseDTO, status.Object)
  GetCustomerVouchers(ctx context.Context, vouchersQuery *query.GetCustomerVouchersQuery) ([]dto.VoucherResponseDTO, status.Object)
}

func NewCustomerQuery(config CustomerQueryConfig) ICustomerQuery {
  return &customerQueryService{
    CustomerQueryConfig: config,
  }
}

func DefaultCustomerQueryConfig(repo repository.ICustomer) CustomerQueryConfig {
  return CustomerQueryConfig{
    GetById:      query.NewGetCustomerByIdHandler(repo),
    GetAddresses: query.NewGetCustomerAddresses(repo),
    GetVouchers:  query.NewGetCustomerVouchers(repo),
  }
}

type CustomerQueryConfig struct {
  GetById      query.IGetCustomerByIdsHandler
  GetAddresses query.IGetCustomerAddressesHandler
  GetVouchers  query.IGetCustomerVouchersHandler
}

type customerQueryService struct {
  CustomerQueryConfig
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
