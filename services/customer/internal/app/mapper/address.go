package mapper

import (
  "rashop/services/customer/internal/app/dto"
  "rashop/services/customer/internal/domain/entity"
)

func ToAddressResponseDTO(address *entity.Address) dto.AddressResponseDTO {
  return dto.AddressResponseDTO{
    Id:             address.Id,
    StreetAddress1: address.StreetAddress1,
    StreetAddress2: address.StreetAddress2,
    City:           address.City,
    State:          address.State,
    PostalCode:     address.PostalCode,
    LastModifiedAt: address.LastModifiedAt,
    CreatedAt:      address.CreatedAt,
  }
}
