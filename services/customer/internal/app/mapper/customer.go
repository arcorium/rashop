package mapper

import (
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "rashop/services/customer/internal/app/dto"
  "rashop/services/customer/internal/domain/entity"
)

func ToCustomerResponseDTO(customer *entity.Customer) dto.CustomerResponseDTO {
  return dto.CustomerResponseDTO{
    Id:             customer.Id,
    Username:       customer.Name.User,
    FirstName:      customer.Name.First,
    LastName:       customer.Name.Last,
    Email:          customer.Email,
    Balance:        customer.Balance.Total,
    Point:          customer.Balance.Point,
    IsVerified:     customer.IsVerified,
    IsDisabled:     customer.IsDisabled,
    LastModifiedAt: customer.LastModifiedAt,
    CreatedAt:      customer.CreatedAt,
    Addresses:      sharedUtil.CastSliceP(customer.ShippingAddresses.Elements(), ToAddressResponseDTO),
    Vouchers:       sharedUtil.CastSliceP(customer.Vouchers.Elements(), ToVoucherResponseDTO),
  }
}
