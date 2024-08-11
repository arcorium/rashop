package dto

import (
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type VoucherResponseDTO struct {
  Id      types.Id
  AddedAt time.Time
}
