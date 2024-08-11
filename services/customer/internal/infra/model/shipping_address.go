package model

import (
  "github.com/arcorium/rashop/shared/types"
  "github.com/arcorium/rashop/shared/util/repo"
  "github.com/uptrace/bun"
  "mini-shop/services/user/internal/domain/entity"
  "time"
)

type ShippingAddressOption = repo.DataAccessModelMapOption[*entity.Address, *ShippingAddress]

func FromAddressDomain(userId types.Id, ent *entity.Address, options ...ShippingAddressOption) ShippingAddress {
  addr := ShippingAddress{
    Id:             ent.Id.String(),
    UserId:         userId.String(),
    StreetAddress1: ent.StreetAddress1,
    StreetAddress2: ent.StreetAddress2,
    City:           ent.City,
    State:          ent.State,
    PostalCode:     ent.PostalCode,
    UpdatedAt:      ent.LastModifiedAt,
    CreatedAt:      ent.CreatedAt,
  }

  for _, option := range options {
    option(ent, &addr)
  }
  return addr
}

type ShippingAddress struct {
  bun.BaseModel `bun:"table:shipping_addresses,alias:sa"`

  Id             string    `bun:",type:uuid,nullzero,pk"`
  UserId         string    `bun:",type:uuid,nullzero,notnull,pk"`
  StreetAddress1 string    `bun:",nullzero,notnull"`
  StreetAddress2 string    `bun:",nullzero"`
  City           string    `bun:",notnull"`
  State          string    `bun:",nullzero"`
  PostalCode     uint32    `bun:",notnull,nullzero"`
  UpdatedAt      time.Time `bun:",nullzero"`
  CreatedAt      time.Time `bun:",notnull,nullzero"`

  Customer *Customer `bun:"rel:belongs-to,join:user_id=user_id"`
}

func (s *ShippingAddress) ToDomain() (entity.Address, error) {
  id, err := types.IdFromString(s.Id)
  if err != nil {
    return entity.Address{}, err
  }

  return entity.Address{
    Id:             id,
    StreetAddress1: s.StreetAddress1,
    StreetAddress2: s.StreetAddress2,
    City:           s.City,
    State:          s.State,
    PostalCode:     s.PostalCode,
    CreatedAt:      s.CreatedAt,
    LastModifiedAt: s.UpdatedAt,
  }, nil
}
