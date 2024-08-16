package command

import (
  "github.com/arcorium/rashop/shared/types"
  "rashop/services/customer/internal/domain/entity"
  "time"
)

type CreateCustomerCommand struct {
  Username  string
  FirstName string
  LastName  types.NullableString
  Email     types.Email
  Password  types.Password
}

func (c *CreateCustomerCommand) ToDomain() (entity.Customer, error) {
  id, err := types.NewId()
  if err != nil {
    return entity.Customer{}, err
  }

  cust := entity.Customer{
    Id:        id,
    CreatedAt: time.Now(),
  }

  hashedPassword, err := c.Password.Hash()
  if err != nil {
    return cust, err
  }

  cust.Name.User = c.Username
  cust.Name.First = c.FirstName
  cust.Email = c.Email
  cust.Password = hashedPassword
  types.SetOnNonNull(&cust.Name.Last, c.LastName)

  return cust, nil
}

type UpdateCustomerCommand struct {
  CustomerId types.Id
  Username   types.NullableString
  FirstName  types.NullableString
  LastName   types.NullableString
  Email      types.NullableEmail
}

func (c *UpdateCustomerCommand) ToPredicate() types.UnaryPredicate[*entity.Customer] {
  return func(customer *entity.Customer) {
    types.SetOnNonNull(&customer.Name.User, c.Username)
    types.SetOnNonNull(&customer.Name.First, c.FirstName)
    types.SetOnNonNull(&customer.Name.Last, c.LastName)
    types.SetOnNonNull(&customer.Email, c.Email)
  }
}

type DisableCustomerCommand struct {
  CustomerId types.Id
}

type EnableCustomerCommand struct {
  CustomerId types.Id
}

type ResetCustomerPasswordCommand struct {
  Token       string
  LogoutAll   bool
  NewPassword types.Password
}

type UpdateBalanceOperator string

const (
  OperatorSet UpdateBalanceOperator = "set"
  OperatorMod UpdateBalanceOperator = "mod" // Modifier (increase/subtract)
)

type UpdateCustomerBalanceCommand struct {
  CustomerId types.Id
  Operator   UpdateBalanceOperator
  Balance    int64
  Point      int64
}

type UpdateCustomerPhotoCommand struct {
  CustomerId types.Id
  MediaId    types.Id
}

type UpdateCustomerPasswordCommand struct {
  CustomerId   types.Id
  LastPassword types.Password
  NewPassword  types.Password
}

type VerifyCustomerEmailCommand struct {
  Token string
}

type ForgotCustomerPasswordRequestCommand struct {
  Email types.Email
}

type VerificationCustomerEmailRequestCommand struct {
  CustomerId types.Id
}

type AddCustomerAddressCommand struct {
  CustomerId     types.Id
  StreetAddress1 string `validate:"required"`
  StreetAddress2 types.NullableString
  City           string `validate:"required"`
  State          string `validate:"required"`
  PostalCode     uint32 `validate:"required"`
}

func (a *AddCustomerAddressCommand) ToDomain() (entity.Address, error) {
  id, err := types.NewId()
  if err != nil {
    return entity.Address{}, err
  }
  return entity.Address{
    Id:             id,
    StreetAddress1: a.StreetAddress1,
    StreetAddress2: a.StreetAddress2.ValueOr(""),
    City:           a.City,
    State:          a.State,
    PostalCode:     a.PostalCode,
    CreatedAt:      time.Now(),
  }, nil
}

type UpdateCustomerAddressCommand struct {
  CustomerId     types.Id
  AddressId      types.Id
  StreetAddress1 types.NullableString
  StreetAddress2 types.NullableString
  City           types.NullableString
  State          types.NullableString
  PostalCode     types.NullableUInt32
}

func (a *UpdateCustomerAddressCommand) ToPredicate() types.UnaryPredicate[*entity.Address] {
  return func(address *entity.Address) {
    types.SetOnNonNull(&address.StreetAddress1, a.StreetAddress1)
    types.SetOnNonNull(&address.StreetAddress2, a.StreetAddress2)
    types.SetOnNonNull(&address.City, a.City)
    types.SetOnNonNull(&address.State, a.State)
    types.SetOnNonNull(&address.PostalCode, a.PostalCode)
  }
}

type DeleteCustomerAddressCommand struct {
  CustomerId types.Id
  AddressId  types.Id
}

type AddCustomerVoucherCommand struct {
  CustomerId types.Id
  VoucherId  types.Id
}

func (c *AddCustomerVoucherCommand) ToDomain() entity.Voucher {
  return entity.Voucher{
    Id:          c.VoucherId,
    IsBeingUsed: false,
    CreatedAt:   time.Now(),
  }
}

type UpdateCustomerVoucherCommand struct {
  CustomerId  types.Id
  VoucherId   types.Id
  IsBeingUsed bool
}

func (a *UpdateCustomerVoucherCommand) ToPredicate() types.UnaryPredicate[*entity.Voucher] {
  return func(voucher *entity.Voucher) {
    voucher.IsBeingUsed = a.IsBeingUsed
  }
}

type DeleteCustomerVoucherCommand struct {
  CustomerId types.Id
  VoucherId  types.Id
}

type SetCustomerDefaultAddressCommand struct {
  CustomerId types.Id
  AddressId  types.Id
}
