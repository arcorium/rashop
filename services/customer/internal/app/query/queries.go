package query

import (
  sharedDto "github.com/arcorium/rashop/shared/dto"
  "github.com/arcorium/rashop/shared/types"
)

type GetCustomersQuery struct {
  sharedDto.PagedElementDTO
}

type GetCustomerByIdsQuery struct {
  CustomerIds []types.Id
}

type GetCustomerAddressesQuery struct {
  CustomerId types.Id
}

type GetCustomerVouchersQuery struct {
  CustomerId types.Id
}
