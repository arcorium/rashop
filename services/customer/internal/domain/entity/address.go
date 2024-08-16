package entity

import (
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type Address struct {
  Id             types.Id
  StreetAddress1 string
  StreetAddress2 string
  City           string
  State          string
  PostalCode     uint32
  CreatedAt      time.Time
  LastModifiedAt time.Time
}

type preservedAddressFields struct {
  Id        types.Id
  CreatedAt time.Time
}

func (a *Address) PreserveFields() preservedAddressFields {
  return preservedAddressFields{
    Id:        a.Id,
    CreatedAt: a.CreatedAt,
  }
}

func (a *Address) RestorePreserved(fields *preservedAddressFields) {
  a.Id = fields.Id
  a.CreatedAt = fields.CreatedAt
}

func ReorderAddresses(defaultId string, addresses []Address) []Address {
  // Put default address in front
  if len(defaultId) == 0 || len(addresses) == 0 {
    return addresses
  }
  // Create new addresses without the default
  var defaultsIdx int
  var temp []Address // Without default shipping address
  for i, v := range addresses {
    if v.Id.EqWithString(defaultId) {
      defaultsIdx = i
      continue
    }
    temp = append(temp, v)
  }
  // Prepend it
  return append([]Address{addresses[defaultsIdx]}, temp...)
}
