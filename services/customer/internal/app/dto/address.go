package dto

import (
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type AddressResponseDTO struct {
  Id             types.Id
  StreetAddress1 string
  StreetAddress2 string
  City           string
  State          string
  PostalCode     uint32
  LastModifiedAt time.Time
  CreatedAt      time.Time
}
