package event

import (
  "github.com/arcorium/rashop/shared/types"
  "time"
)

var _ types.Event = (*CustomerCreatedV1)(nil)

type CustomerCreatedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
  Email      types.Email
  Password   types.HashedPassword // Hashed
  Username   string
  FirstName  string
  LastName   string
  Balance    uint64
  Point      uint64
  IsVerified bool
  IsDisabled bool
  CreatedAt  time.Time
}

func (c *CustomerCreatedV1) EventName() string {
  return CustomerCreatedEvent
}

func (c *CustomerCreatedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerUpdatedV1)(nil)

type CustomerUpdatedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
  Username   string
  FirstName  string
  LastName   string
  Email      types.Email
}

func (c *CustomerUpdatedV1) EventName() string {
  return CustomerUpdatedEvent
}

func (c *CustomerUpdatedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerPasswordUpdatedV1)(nil)

type CustomerPasswordUpdatedV1 struct {
  types.DomainEventBaseV1
  CustomerId  string
  NewPassword types.HashedPassword // hashed
}

func (c *CustomerPasswordUpdatedV1) EventName() string {
  return CustomerPasswordUpdatedEvent
}

func (c *CustomerPasswordUpdatedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerDeletedV1)(nil)

type CustomerDeletedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
}

func (c *CustomerDeletedV1) EventName() string {
  return CustomerDeletedEvent
}

func (c *CustomerDeletedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerBalanceUpdatedV1)(nil)

type CustomerBalanceUpdatedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
  Balance    uint64
  Point      uint64
}

func (c *CustomerBalanceUpdatedV1) EventName() string {
  return CustomerBalanceUpdatedEvent
}

func (c *CustomerBalanceUpdatedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerStatusUpdatedV1)(nil)

type CustomerStatusUpdatedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
  Status     bool
}

func (c *CustomerStatusUpdatedV1) EventName() string {
  return CustomerStatusUpdatedEvent
}

func (c *CustomerStatusUpdatedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerPhotoUpdatedV1)(nil)

type CustomerPhotoUpdatedV1 struct {
  types.DomainEventBaseV1
  CustomerId  string
  LastMediaId string // Used when user set new photo
  NewMediaId  types.Id
}

func (c *CustomerPhotoUpdatedV1) EventName() string {
  return CustomerPhotoUpdatedEvent
}

func (c *CustomerPhotoUpdatedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerEmailVerifiedV1)(nil)

type CustomerEmailVerifiedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
}

func (c *CustomerEmailVerifiedV1) EventName() string {
  return CustomerEmailVerifiedEvent
}

func (c *CustomerEmailVerifiedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerForgotPasswordRequestedV1)(nil)

type CustomerForgotPasswordRequestedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
  Email      string
  Username   string
}

func (c *CustomerForgotPasswordRequestedV1) EventName() string {
  return CustomerForgotPasswordRequestedEvent
}

func (c *CustomerForgotPasswordRequestedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerEmailVerificationRequestedV1)(nil)

type CustomerEmailVerificationRequestedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
  Email      string
  Username   string
}

func (c *CustomerEmailVerificationRequestedV1) EventName() string {
  return CustomerEmailVerificationRequestedEvent
}

func (c *CustomerEmailVerificationRequestedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerAddressAddedV1)(nil)

type CustomerAddressAddedV1 struct {
  types.DomainEventBaseV1
  CustomerId     string
  AddressId      types.Id
  StreetAddress1 string
  StreetAddress2 string
  City           string
  State          string
  PostalCode     uint32
}

func (c *CustomerAddressAddedV1) EventName() string {
  return CustomerAddressAddedEvent
}

func (c *CustomerAddressAddedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerAddressDeletedV1)(nil)

type CustomerAddressDeletedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
  AddressId  types.Id
}

func (c *CustomerAddressDeletedV1) EventName() string {
  return CustomerAddressesDeletedEvent
}

func (c *CustomerAddressDeletedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerAddressUpdatedV1)(nil)

type CustomerAddressUpdatedV1 struct {
  types.DomainEventBaseV1
  CustomerId     string
  AddressId      types.Id
  StreetAddress1 string
  StreetAddress2 string
  City           string
  State          string
  PostalCode     uint32
}

func (c *CustomerAddressUpdatedV1) EventName() string {
  return CustomerAddressUpdatedEvent
}

func (c *CustomerAddressUpdatedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerVoucherAddedV1)(nil)

type CustomerVoucherAddedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
  VoucherId  types.Id
}

func (c *CustomerVoucherAddedV1) EventName() string {
  return CustomerVouchersAddedEvent
}

func (c *CustomerVoucherAddedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerVoucherDeletedV1)(nil)

type CustomerVoucherDeletedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
  VoucherId  types.Id
}

func (c *CustomerVoucherDeletedV1) EventName() string {
  return CustomerVouchersDeletedEvent
}

func (c *CustomerVoucherDeletedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerVoucherUpdatedV1)(nil)

type CustomerVoucherUpdatedV1 struct {
  types.DomainEventBaseV1
  CustomerId  string
  VoucherId   types.Id
  IsBeingUsed bool
}

func (c *CustomerVoucherUpdatedV1) EventName() string {
  return CustomerVoucherUpdatedEvent
}

func (c *CustomerVoucherUpdatedV1) Key() string {
  return c.CustomerId
}

var _ types.Event = (*CustomerDefaultAddressUpdatedV1)(nil)

type CustomerDefaultAddressUpdatedV1 struct {
  types.DomainEventBaseV1
  CustomerId string
  AddressId  types.Id
}

func (c *CustomerDefaultAddressUpdatedV1) EventName() string {
  return CustomerDefaultAddressUpdatedEvent
}

func (c *CustomerDefaultAddressUpdatedV1) Key() string {
  return c.CustomerId
}
