package model

import (
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "github.com/arcorium/rashop/shared/util/repo"
  "github.com/uptrace/bun"
  "mini-shop/services/user/internal/domain/entity"
  vob "mini-shop/services/user/internal/domain/valueobject"
)

type CustomerOption = repo.DataAccessModelMapOption[*entity.Customer, *Customer]

func FromCustomerDomain(ent *entity.Customer, options ...CustomerOption) Customer {
  var saId string
  if sa := ent.DefaultShippingAddress(); sa != nil {
    saId = sa.Id.String()
  }

  cust := Customer{
    UserId:                   ent.Id.String(),
    FirstName:                ent.Name.First,
    LastName:                 ent.Name.Last,
    PhotoId:                  ent.PhotoId.String(),
    Balance:                  ent.Balance.Total,
    Point:                    ent.Balance.Point,
    DefaultShippingAddressId: saId,
    IsDisabled:               ent.IsDisabled,
    User: &User{
      Id:         ent.Id.String(),
      Username:   ent.Name.User,
      Email:      ent.Email.String(),
      Password:   ent.Password.String(),
      IsVerified: ent.IsVerified,
      UpdatedAt:  ent.LastModifiedAt,
      CreatedAt:  ent.CreatedAt,
    },
    Vouchers: sharedUtil.CastSliceP(ent.Vouchers.Elements(), func(vc *entity.Voucher) Voucher {
      return FromVoucherDomain(ent.Id, vc)
    }),
    ShippingAddresses: sharedUtil.CastSliceP(ent.ShippingAddresses.Elements(), func(addr *entity.Address) ShippingAddress {
      return FromAddressDomain(ent.Id, addr)
    }),
  }

  for _, opt := range options {
    opt(ent, &cust)
  }
  return cust
}

type Customer struct {
  bun.BaseModel `bun:"table:customers,alias:c"`

  UserId                   string `bun:",type:uuid,nullzero,pk"`
  FirstName                string `bun:",nullzero,notnull"`
  LastName                 string `bun:","`
  PhotoId                  string `bun:",type:uuid"`
  Balance                  uint64 `bun:",default:0"`
  Point                    uint64 `bun:",default:0"`
  DefaultShippingAddressId string `bun:"default_sa_id,type:uuid,nullzero"`
  IsDisabled               bool   `bun:",default:false"`

  User              *User             `bun:"rel:belongs-to,join:user_id=id"`
  Vouchers          []Voucher         `bun:"rel:has-many,join:user_id=user_id"`
  ShippingAddresses []ShippingAddress `bun:"rel:has-many,join:user_id=user_id"`
}

func (c *Customer) ToDomain() (entity.Customer, error) {
  id, err := types.IdFromString(c.UserId)
  if err != nil {
    return entity.Customer{}, err
  }

  email, err := types.EmailFromString(c.User.Email)
  if err != nil {
    return entity.Customer{}, err
  }

  photoId, err := types.IdFromString(c.PhotoId)
  if err != nil {
    return entity.Customer{}, err
  }

  addresses, ierr := sharedUtil.CastSliceErrsP(c.ShippingAddresses, repo.ToDomainErr[*ShippingAddress, entity.Address])
  if !ierr.IsNil() {
    return entity.Customer{}, ierr
  }

  // Put default address in front
  addresses = entity.ReorderAddresses(c.DefaultShippingAddressId, addresses)

  vouchers, ierr := sharedUtil.CastSliceErrsP(c.Vouchers, repo.ToDomainErr[*Voucher, entity.Voucher])
  if !ierr.IsNil() {
    return entity.Customer{}, nil
  }

  return entity.Customer{
    Id: id,
    Name: vob.Name{
      User:  c.User.Username,
      First: c.FirstName,
      Last:  c.LastName,
    },
    Email:    email,
    Password: types.HashedPassword(c.User.Password),
    PhotoId:  photoId,
    Balance: vob.Balance{
      Total: c.Balance,
      Point: c.Point,
    },
    IsVerified:        c.User.IsVerified,
    IsDisabled:        c.IsDisabled,
    LastModifiedAt:    c.User.UpdatedAt,
    CreatedAt:         c.User.CreatedAt,
    ShippingAddresses: types.NewChildEntityHelper[types.Id](addresses),
    Vouchers:          types.NewChildEntityHelper[types.Id](vouchers),
  }, nil
}
