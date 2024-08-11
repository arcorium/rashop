package entity

import (
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type Voucher struct {
  Id          types.Id
  IsBeingUsed bool // If the voucher currently used for order that is not completed yet
  CreatedAt   time.Time
}

type preservedVoucherFields struct {
  Id        types.Id
  CreatedAt time.Time
}

func (a *Voucher) PreserveFields() preservedVoucherFields {
  return preservedVoucherFields{
    Id:        a.Id,
    CreatedAt: a.CreatedAt,
  }
}

func (a *Voucher) RestorePreserved(fields *preservedVoucherFields) {
  a.Id = fields.Id
  a.CreatedAt = fields.CreatedAt
}
