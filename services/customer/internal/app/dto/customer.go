package dto

import (
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type CustomerResponseDTO struct {
  Id             types.Id
  Username       string
  FirstName      string
  LastName       string
  Email          types.Email
  Balance        uint64
  Point          uint64
  IsVerified     bool
  IsDisabled     bool
  LastModifiedAt time.Time
  CreatedAt      time.Time

  Addresses []AddressResponseDTO
  Vouchers  []VoucherResponseDTO
}
