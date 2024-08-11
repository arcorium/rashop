package model

import (
  "github.com/arcorium/rashop/shared/types"
  "github.com/arcorium/rashop/shared/util/repo"
  "github.com/uptrace/bun"
  "mini-shop/services/user/internal/domain/entity"
  "time"
)

type VoucherOption = repo.DataAccessModelMapOption[*entity.Voucher, *Voucher]

func FromVoucherDomain(userId types.Id, ent *entity.Voucher, options ...VoucherOption) Voucher {
  vcr := Voucher{
    UserId:      userId.String(),
    VoucherId:   ent.Id.String(),
    IsBeingUsed: ent.IsBeingUsed,
    CreatedAt:   ent.CreatedAt,
  }

  for _, option := range options {
    option(ent, &vcr)
  }
  return vcr
}

type Voucher struct {
  bun.BaseModel `bun:"table:vouchers,alias:v"`
  UserId        string    `bun:",type:uuid,nullzero,pk"`
  VoucherId     string    `bun:",type:uuid,nullzero,pk"`
  IsBeingUsed   bool      `bun:",default:false"`
  CreatedAt     time.Time `bun:",nullzero,notnull"`

  Customer *Customer `bun:"rel:belongs-to,join:user_id=user_id"`
}

func (v *Voucher) ToDomain() (entity.Voucher, error) {
  id, err := types.IdFromString(v.VoucherId)
  if err != nil {
    return entity.Voucher{}, err
  }

  return entity.Voucher{
    Id:          id,
    IsBeingUsed: v.IsBeingUsed,
    CreatedAt:   v.CreatedAt,
  }, nil
}
