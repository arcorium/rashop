package query

import "github.com/arcorium/rashop/shared/types"

type GetCustomerByIdsQuery struct {
  CustomerIds []types.Id
}

type GetCustomerAddressesQuery struct {
  CustomerId types.Id
}

type GetCustomerVouchersQuery struct {
  CustomerId types.Id
}
