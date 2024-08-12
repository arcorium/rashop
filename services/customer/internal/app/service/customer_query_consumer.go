package service

import (
  "context"
  "github.com/arcorium/rashop/shared/status"
  "mini-shop/services/user/internal/app/query/consumer"
  "mini-shop/services/user/internal/domain/event"
  "mini-shop/services/user/internal/domain/repository"
)

type ICustomerQueryConsumer interface {
  AddressAddedV1(ctx context.Context, ev *event.CustomerAddressAddedV1) status.Object
  AddressDeletedV1(ctx context.Context, ev *event.CustomerAddressDeletedV1) status.Object
  AddressUpdatedV1(ctx context.Context, ev *event.CustomerAddressUpdatedV1) status.Object
  BalanceUpdatedV1(ctx context.Context, ev *event.CustomerBalanceUpdatedV1) status.Object
  CreatedV1(ctx context.Context, ev *event.CustomerCreatedV1) status.Object
  DefaultAddressSetV1(ctx context.Context, ev *event.CustomerDefaultAddressUpdatedV1) status.Object
  EmailVerifiedV1(ctx context.Context, ev *event.CustomerEmailVerifiedV1) status.Object
  PasswordUpdatedV1(ctx context.Context, ev *event.CustomerPasswordUpdatedV1) status.Object
  PhotoUpdatedV1(ctx context.Context, ev *event.CustomerPhotoUpdatedV1) status.Object
  StatusUpdatedV1(ctx context.Context, ev *event.CustomerStatusUpdatedV1) status.Object
  UpdatedV1(ctx context.Context, ev *event.CustomerUpdatedV1) status.Object
  VoucherAddedV1(ctx context.Context, ev *event.CustomerVoucherAddedV1) status.Object
  VoucherDeletedV1(ctx context.Context, ev *event.CustomerVoucherDeletedV1) status.Object
  VoucherUpdatedV1(ctx context.Context, ev *event.CustomerVoucherUpdatedV1) status.Object
}

func NewCustomerQueryConsumer(config CustomerQueryConsumerConfig) ICustomerQueryConsumer {
  return &customerQueryConsumerService{
    CustomerQueryConsumerConfig: config,
  }
}

func DefaultCustomerQueryConsumerConfig(repo repository.ICustomer) CustomerQueryConsumerConfig {
  return CustomerQueryConsumerConfig{
    AddressAdded:      consumer.NewCustomerAddressAddedConsumer(repo),
    AddressDeleted:    consumer.NewCustomerAddressDeletedConsumer(repo),
    AddressUpdated:    consumer.NewCustomerAddressUpdatedConsumer(repo),
    BalanceUpdated:    consumer.NewCustomerBalanceUpdatedConsumer(repo),
    Created:           consumer.NewCustomerCreatedConsumer(repo),
    DefaultAddressSet: consumer.NewCustomerDefaultAddressSetConsumer(repo),
    EmailVerified:     consumer.NewCustomerEmailVerifiedConsumer(repo),
    PasswordUpdated:   consumer.NewCustomerPasswordUpdatedConsumer(repo),
    PhotoUpdated:      consumer.NewCustomerPhotoUpdatedConsumer(repo),
    StatusUpdated:     consumer.NewCustomerStatusUpdatedConsumer(repo),
    Updated:           consumer.NewCustomerUpdatedConsumer(repo),
    VouchersAdded:     consumer.NewCustomerVoucherAddedConsumer(repo),
    VoucherDeleted:    consumer.NewCustomerVoucherDeletedConsumer(repo),
    VoucherUpdated:    consumer.NewCustomerVoucherUpdatedConsumer(repo),
  }
}

type CustomerQueryConsumerConfig struct {
  AddressAdded      consumer.ICustomerAddressAddedConsumer
  AddressDeleted    consumer.ICustomerAddressDeletedConsumer
  AddressUpdated    consumer.ICustomerAddressUpdatedConsumer
  BalanceUpdated    consumer.ICustomerBalanceUpdatedConsumer
  Created           consumer.ICustomerCreatedConsumer
  DefaultAddressSet consumer.ICustomerDefaultAddressSetConsumer
  EmailVerified     consumer.ICustomerEmailVerifiedConsumer
  PasswordUpdated   consumer.ICustomerPasswordUpdatedConsumer
  PhotoUpdated      consumer.ICustomerPhotoUpdatedConsumer
  StatusUpdated     consumer.ICustomerStatusUpdatedConsumer
  Updated           consumer.ICustomerUpdatedConsumer
  VouchersAdded     consumer.ICustomerVoucherAddedConsumer
  VoucherDeleted    consumer.ICustomerVoucherDeletedConsumer
  VoucherUpdated    consumer.ICustomerVoucherUpdatedConsumer
}

type customerQueryConsumerService struct {
  CustomerQueryConsumerConfig
}

func (c *customerQueryConsumerService) AddressAddedV1(ctx context.Context, ev *event.CustomerAddressAddedV1) status.Object {
  return c.AddressAdded.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) AddressDeletedV1(ctx context.Context, ev *event.CustomerAddressDeletedV1) status.Object {
  return c.AddressDeleted.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) AddressUpdatedV1(ctx context.Context, ev *event.CustomerAddressUpdatedV1) status.Object {
  return c.AddressUpdated.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) BalanceUpdatedV1(ctx context.Context, ev *event.CustomerBalanceUpdatedV1) status.Object {
  return c.BalanceUpdated.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) CreatedV1(ctx context.Context, ev *event.CustomerCreatedV1) status.Object {
  return c.Created.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) DefaultAddressSetV1(ctx context.Context, ev *event.CustomerDefaultAddressUpdatedV1) status.Object {
  return c.DefaultAddressSet.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) EmailVerifiedV1(ctx context.Context, ev *event.CustomerEmailVerifiedV1) status.Object {
  return c.EmailVerified.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) PasswordUpdatedV1(ctx context.Context, ev *event.CustomerPasswordUpdatedV1) status.Object {
  return c.PasswordUpdated.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) PhotoUpdatedV1(ctx context.Context, ev *event.CustomerPhotoUpdatedV1) status.Object {
  return c.PhotoUpdated.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) StatusUpdatedV1(ctx context.Context, ev *event.CustomerStatusUpdatedV1) status.Object {
  return c.StatusUpdated.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) UpdatedV1(ctx context.Context, ev *event.CustomerUpdatedV1) status.Object {
  return c.Updated.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) VoucherAddedV1(ctx context.Context, ev *event.CustomerVoucherAddedV1) status.Object {
  return c.VouchersAdded.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) VoucherDeletedV1(ctx context.Context, ev *event.CustomerVoucherDeletedV1) status.Object {
  return c.VoucherDeleted.Handle(ctx, ev)
}

func (c *customerQueryConsumerService) VoucherUpdatedV1(ctx context.Context, ev *event.CustomerVoucherUpdatedV1) status.Object {
  return c.VoucherUpdated.Handle(ctx, ev)
}
