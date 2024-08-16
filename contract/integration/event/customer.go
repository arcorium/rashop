package intev

import (
  "github.com/arcorium/rashop/shared/types"
)

const (
  CustomerCreatedEvent                    = "customer.created"
  CustomerDisabledEvent                   = "customer.disabled"
  CustomerEnabledEvent                    = "customer.enabled"
  CustomerEmailUpdatedEvent               = "customer.email.updated"
  CustomerPhotoUpdatedEvent               = "customer.photo.updated"
  CustomerResetPasswordRequestedEvent     = "customer.reset-password.requested"
  CustomerEmailVerificationRequestedEvent = "customer.email-verification.requested"
)

var _ types.Event = (*CustomerCreatedV1)(nil)

func NewCustomerCreatedV1(customerId types.Id, email types.Email, username string) *CustomerCreatedV1 {
  return &CustomerCreatedV1{
    IntegrationV1: NewV1(),
    CustomerId:    customerId.String(),
    Email:         email.String(),
    Username:      username,
  }
}

type CustomerCreatedV1 struct {
  IntegrationV1
  CustomerId string
  Email      string
  Username   string
}

func (c *CustomerCreatedV1) EventName() string {
  return CustomerCreatedEvent
}

func (c *CustomerCreatedV1) Key() (string, bool) {
  return c.CustomerId, true
}

var _ types.Event = (*CustomerEmailChangedV1)(nil)

func NewCustomerEmailChangedV1(customerId types.Id, newEmail types.Email) *CustomerEmailChangedV1 {
  return &CustomerEmailChangedV1{
    IntegrationV1: NewV1(),
    CustomerId:    customerId.String(),
    NewEmail:      newEmail.String(),
  }
}

type CustomerEmailChangedV1 struct {
  IntegrationV1
  CustomerId string
  NewEmail   string
}

func (c *CustomerEmailChangedV1) EventName() string {
  return CustomerEmailUpdatedEvent
}

func (c *CustomerEmailChangedV1) Key() (string, bool) {
  return c.CustomerId, true
}

var _ types.Event = (*CustomerPhotoChangedV1)(nil)

func NewCustomerPhotoChangedV1(customerId, lastMediaId, newMediaId types.Id) *CustomerPhotoChangedV1 {
  return &CustomerPhotoChangedV1{
    IntegrationV1: NewV1(),
    CustomerId:    customerId.String(),
    LastMediaId:   lastMediaId.String(),
    NewMediaId:    newMediaId.String(),
  }
}

type CustomerPhotoChangedV1 struct {
  IntegrationV1
  CustomerId  string
  LastMediaId string
  NewMediaId  string
}

func (c *CustomerPhotoChangedV1) EventName() string {
  return CustomerPhotoUpdatedEvent
}

func (c *CustomerPhotoChangedV1) Key() (string, bool) {
  return c.CustomerId, true
}

var _ types.Event = (*CustomerResetPasswordRequestedV1)(nil)

func NewCustomerResetPasswordRequestedV1(customerId types.Id, email types.Email, username string) *CustomerResetPasswordRequestedV1 {
  return &CustomerResetPasswordRequestedV1{
    IntegrationV1: NewV1(),
    CustomerId:    customerId.String(),
    Email:         email.String(),
    Username:      username,
  }
}

type CustomerResetPasswordRequestedV1 struct {
  IntegrationV1
  CustomerId string
  Email      string
  Username   string
}

func (c *CustomerResetPasswordRequestedV1) EventName() string {
  return CustomerResetPasswordRequestedEvent
}

func (c *CustomerResetPasswordRequestedV1) Key() (string, bool) {
  return c.CustomerId, true
}

var _ types.Event = (*CustomerEmailVerificationRequestedV1)(nil)

func NewCustomerEmailVerificationRequestedV1(customerId types.Id, email types.Email, username string) *CustomerEmailVerificationRequestedV1 {
  return &CustomerEmailVerificationRequestedV1{
    IntegrationV1: NewV1(),
    CustomerId:    customerId.String(),
    Email:         email.String(),
    Username:      username,
  }
}

type CustomerEmailVerificationRequestedV1 struct {
  IntegrationV1
  CustomerId string
  Email      string
  Username   string
}

func (c *CustomerEmailVerificationRequestedV1) EventName() string {
  return CustomerEmailVerificationRequestedEvent
}

func (c *CustomerEmailVerificationRequestedV1) Key() (string, bool) {
  return c.CustomerId, true
}
