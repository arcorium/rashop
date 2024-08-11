package entity

import (
  "errors"
  "fmt"
  algo "github.com/arcorium/rashop/shared/algorithm"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "mini-shop/services/user/internal/domain/event"
  vob "mini-shop/services/user/internal/domain/valueobject"
  "time"
)

var (
  ErrDisableDisabledCustomer = errors.New("couldn't disabling customer that already disabled")
  ErrEnableEnabledCustomer   = errors.New("couldn't enabling customer that already enabled")
  ErrSamePassword            = errors.New("couldn't reset password with the same old password")
  ErrDifferentPassword       = errors.New("last password doesn't match the current password")
  ErrUserAlreadyVerified     = errors.New("customer already verified")
  ErrUserIsDisabled          = errors.New("couldn't update disabled customer")
  ErrAddressNotFound         = errors.New("customer address not found")
  ErrVoucherNotFound         = errors.New("customer voucher not found")
  ErrVoucherAlreadyExists    = errors.New("customer voucher already exists")
  ErrVoucherIsBeingUsed      = errors.New("couldn't delete being used vouchers")
  ErrNothingToUpdate         = errors.New("nothing to update ")
)

var _ types.Aggregate = (*Customer)(nil)

func CreateCutomer(customer *Customer) (Customer, types.Event) {
  cust := Customer{}
  ev := &event.CustomerCreatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        customer.Id.String(),
    Email:             customer.Email,
    Password:          customer.Password,
    Username:          customer.Name.User,
    FirstName:         customer.Name.First,
    LastName:          customer.Name.Last,
    Balance:           customer.Balance.Total,
    Point:             customer.Balance.Point,
    IsVerified:        customer.IsVerified,
    IsDisabled:        customer.IsDisabled,
    CreatedAt:         customer.CreatedAt,
  }
  cust.AddEvents(ev)
  return cust, ev
}

type Customer struct {
  types.AggregateBase
  types.AggregateHelper

  Id         types.Id
  Name       vob.Name
  Email      types.Email
  Password   types.HashedPassword
  PhotoId    types.Id
  Balance    vob.Balance
  IsVerified bool
  IsDisabled bool

  LastModifiedAt time.Time
  CreatedAt      time.Time

  // Helper
  ShippingAddresses types.ChildEntityHelperWithObject[types.Id, Address] // The first element is the default one
  Vouchers          types.ChildEntityHelperWithObject[types.Id, Voucher]

  emailChanged bool
}

func (c *Customer) preCheck() error {
  if c.IsDisabled {
    return ErrUserIsDisabled
  }
  return nil
}

func (c *Customer) IsEmailChanged() bool {
  return c.emailChanged
}

func (c *Customer) Update(pred types.UnaryPredicate[*Customer]) types.Event {
  cloned := sharedUtil.Clone(c)
  pred(&cloned)

  // Emit event
  ev := &event.CustomerUpdatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    Username:          c.Name.User,
    FirstName:         c.Name.First,
    LastName:          c.Name.Last,
    Email:             c.Email,
  }
  c.AddEvents(ev)
  return ev
}

func (c *Customer) Enable() (types.Event, error) {
  if !c.IsDisabled {
    return nil, ErrEnableEnabledCustomer
  }

  ev := &event.CustomerStatusUpdatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    Status:            c.IsDisabled,
  }
  c.AddEvents(ev)

  return ev, nil
}

func (c *Customer) Disable() (types.Event, error) {
  if c.IsDisabled {
    return nil, ErrDisableDisabledCustomer
  }

  ev := &event.CustomerStatusUpdatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    Status:            true,
  }
  c.AddEvents(ev)

  return ev, nil
}

func (c *Customer) ResetPassword(newPassword types.Password) (types.Event, error) {
  // Check if equal
  if c.Password.Eq(newPassword) {
    return nil, ErrSamePassword
  }

  hashed, err := newPassword.Hash()
  if err != nil {
    return nil, err
  }

  ev := &event.CustomerPasswordUpdatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    NewPassword:       hashed,
  }
  c.AddEvents(ev)

  return ev, nil
}

func (c *Customer) SetBalance(balance, point uint64) types.Event {
  ev := &event.CustomerBalanceUpdatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    Balance:           balance,
    Point:             point,
  }
  c.AddEvents(ev)
  return ev
}

func (c *Customer) ModifyBalance(balanceModifier, pointModifier int64) types.Event {
  balance := c.Balance.Total + uint64(balanceModifier)
  point := c.Balance.Point + uint64(pointModifier)
  return c.SetBalance(balance, point)
}

func (c *Customer) UpdatePhoto(mediaId types.Id) types.Event {
  ev := &event.CustomerPhotoUpdatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    LastMediaId:       c.PhotoId.String(),
    NewMediaId:        mediaId,
  }
  c.AddEvents(ev)
  return ev
}

func (c *Customer) UpdatePassword(lastPassword, newPassword types.Password) (types.Event, error) {
  if !c.Password.Eq(lastPassword) {
    return nil, ErrDifferentPassword
  }
  return c.ResetPassword(newPassword)
}

func (c *Customer) SetDefaultAddress(addressId types.Id) (types.Event, error) {
  if c.ShippingAddresses.Elm[0].Id.Eq(addressId) {
    return nil, ErrNothingToUpdate
  }

  _, ok := c.getAddressIndex(addressId)
  if !ok {
    return nil, ErrAddressNotFound
  }

  ev := &event.CustomerDefaultAddressUpdatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    AddressId:         addressId,
  }
  c.AddEvents(ev)

  return ev, nil
}

func (c *Customer) getAddressIndex(addressId types.Id) (int, bool) {
  // Check if exists
  return algo.IndexOfFunc(c.ShippingAddresses.Elm, addressId, func(address *Address, id types.Id) bool {
    return address.Id == id
  })
}

func (c *Customer) AddAddress(address *Address) types.Event {
  ev := &event.CustomerAddressAddedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    AddressId:         address.Id,
    StreetAddress1:    address.StreetAddress1,
    StreetAddress2:    address.StreetAddress2,
    City:              address.City,
    State:             address.State,
    PostalCode:        address.PostalCode,
  }
  c.AddEvents(ev)
  return ev
}

func (c *Customer) UpdateAddress(addressId types.Id, pred types.UnaryPredicate[*Address]) (types.Event, error) {
  idx, ok := c.getAddressIndex(addressId)
  if !ok {
    return nil, ErrAddressNotFound
  }
  current := &c.ShippingAddresses.Elm[idx]

  cloned := sharedUtil.Clone(current)
  pred(&cloned)

  // Emit event
  ev := &event.CustomerAddressUpdatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    AddressId:         current.Id,
    StreetAddress1:    cloned.StreetAddress1,
    StreetAddress2:    cloned.StreetAddress2,
    City:              cloned.City,
    State:             cloned.State,
    PostalCode:        cloned.PostalCode,
  }
  c.AddEvents(ev)

  return ev, nil
}

func (c *Customer) DeleteAddress(addressId types.Id) (types.Event, error) {
  _, ok := c.getAddressIndex(addressId)
  if !ok {
    return nil, ErrAddressNotFound
  }

  ev := &event.CustomerAddressDeletedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    AddressId:         addressId,
  }
  c.AddEvents(ev)
  return ev, nil
}

func (c *Customer) getVoucherIndex(voucherId types.Id) (int, bool) {
  // Check if exists
  return algo.IndexOfFunc(c.Vouchers.Elm, voucherId, func(voucher *Voucher, id types.Id) bool {
    return voucher.Id == id
  })
}

func (c *Customer) AddVoucher(voucher *Voucher) (types.Event, error) {
  _, ok := c.getVoucherIndex(voucher.Id)
  if ok {
    return nil, ErrVoucherAlreadyExists
  }

  ev := &event.CustomerVoucherAddedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    VoucherId:         voucher.Id,
  }
  c.AddEvents(ev)

  return ev, nil
}

func (c *Customer) DeleteVoucher(voucherId types.Id) (types.Event, error) {
  idx, ok := c.getVoucherIndex(voucherId)
  if !ok {
    return nil, ErrVoucherNotFound
  }

  current := &c.Vouchers.Elm[idx]
  // Prevent delete voucher that is being used
  if current.IsBeingUsed {
    return nil, ErrVoucherIsBeingUsed
  }

  ev := &event.CustomerVoucherDeletedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    VoucherId:         voucherId,
  }
  c.AddEvents(ev)
  return ev, nil
}

func (c *Customer) UpdateVoucher(voucherId types.Id, pred types.UnaryPredicate[*Voucher]) (types.Event, error) {
  idx, ok := c.getVoucherIndex(voucherId)
  if !ok {
    return nil, ErrVoucherNotFound
  }

  current := &c.Vouchers.Elm[idx]
  cloned := sharedUtil.Clone(current)
  pred(&cloned)
  if current.IsBeingUsed == cloned.IsBeingUsed {
    return nil, ErrNothingToUpdate
  }

  ev := &event.CustomerVoucherUpdatedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    VoucherId:         current.Id,
    IsBeingUsed:       current.IsBeingUsed,
  }
  c.AddEvents(ev)
  return ev, nil
}

func (c *Customer) VerifyEmail() (types.Event, error) {
  if err := c.preCheck(); err != nil {
    return nil, err
  }

  if c.IsVerified {
    return nil, ErrUserAlreadyVerified
  }

  ev := &event.CustomerEmailVerifiedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
  }

  c.AddEvents(ev)
  return ev, nil
}

func (c *Customer) EmailVerificationRequest() types.Event {
  ev := &event.CustomerEmailVerificationRequestedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    Email:             c.Email.String(),
    Username:          c.Name.User,
  }
  c.AddEvents(ev)
  return ev
}

func (c *Customer) ForgotPasswordRequest() types.Event {
  ev := &event.CustomerForgotPasswordRequestedV1{
    DomainEventBaseV1: types.NewDomainEventV1(),
    CustomerId:        c.Id.String(),
    Email:             c.Email.String(),
    Username:          c.Name.User,
  }
  c.AddEvents(ev)
  return ev
}

func (c *Customer) DefaultShippingAddress() *Address {
  if len(c.ShippingAddresses.Elm) == 0 {
    return nil
  }
  return &c.ShippingAddresses.Elm[0]
}

func (c *Customer) Identity() string {
  return c.Id.String()
}

func (c *Customer) ApplyEvent(ev types.Event) error {
  switch cur := ev.(type) {
  case *event.CustomerCreatedV1:
    c.Id = types.Must(types.IdFromString(cur.CustomerId))
    c.Email = cur.Email
    c.Password = cur.Password
    c.Name.User = cur.Username
    c.Name.First = cur.FirstName
    c.Name.Last = cur.LastName
    c.IsVerified = cur.IsVerified
    c.IsDisabled = cur.IsDisabled
    c.Balance.Total = cur.Balance
    c.Balance.Point = cur.Point
    c.CreatedAt = cur.CreatedAt
    c.MarkCreated()
  case *event.CustomerUpdatedV1:
    c.Name.User = cur.Username
    c.Name.First = cur.FirstName
    c.Name.Last = cur.LastName
    if c.Email != cur.Email {
      c.Email = cur.Email
      c.emailChanged = true
    }
    c.MarkUpdated()
  case *event.CustomerPasswordUpdatedV1:
    c.Password = cur.NewPassword
    c.MarkUpdated()
  case *event.CustomerBalanceUpdatedV1:
    c.Balance.Total = cur.Balance
    c.Balance.Point = cur.Point
    c.MarkUpdated()
  case *event.CustomerStatusUpdatedV1:
    c.IsDisabled = cur.Status
    c.MarkUpdated()
  case *event.CustomerPhotoUpdatedV1:
    c.PhotoId = cur.NewMediaId
    c.MarkUpdated()
  case *event.CustomerEmailVerifiedV1:
    c.IsVerified = true
    c.MarkUpdated()
  case *event.CustomerAddressAddedV1:
    c.ShippingAddresses.Elm = append(c.ShippingAddresses.Elm, Address{
      Id:             cur.AddressId,
      StreetAddress1: cur.StreetAddress1,
      StreetAddress2: cur.StreetAddress2,
      City:           cur.City,
      State:          cur.State,
      PostalCode:     cur.PostalCode,
      CreatedAt:      cur.OccurredAt(),
    })
    c.ShippingAddresses.Add()
  case *event.CustomerAddressDeletedV1:
    c.ShippingAddresses.Delete(cur.AddressId)
  case *event.CustomerAddressUpdatedV1:
    index, _ := c.getAddressIndex(cur.AddressId)
    c.ShippingAddresses.Elm[index] = Address{
      StreetAddress1: cur.StreetAddress1,
      StreetAddress2: cur.StreetAddress2,
      City:           cur.City,
      State:          cur.State,
      PostalCode:     cur.PostalCode,
      LastModifiedAt: cur.OccurredAt(),
    }
    c.ShippingAddresses.Update(index)
  case *event.CustomerVoucherAddedV1:
    c.Vouchers.Elm = append(c.Vouchers.Elm, Voucher{
      Id:          cur.VoucherId,
      IsBeingUsed: false,
      CreatedAt:   cur.OccurredAt(),
    })
    c.Vouchers.Add()
  case *event.CustomerVoucherDeletedV1:
    c.Vouchers.Delete(cur.VoucherId)
  case *event.CustomerVoucherUpdatedV1:
    index, _ := c.getVoucherIndex(cur.VoucherId)
    c.Vouchers.Elm[index] = Voucher{
      IsBeingUsed: cur.IsBeingUsed,
      CreatedAt:   cur.OccurredAt(),
    }
    c.Vouchers.Update(index)
  case *event.CustomerDefaultAddressUpdatedV1:
    c.ShippingAddresses.Elm = ReorderAddresses(cur.AddressId.String(), c.ShippingAddresses.Elm)
    c.MarkUpdated()
  default:
    return fmt.Errorf("unexpected event type %T", ev)
  }
  return nil
}
