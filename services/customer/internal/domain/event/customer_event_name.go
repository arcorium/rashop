package event

const (
  CustomerCreatedEvent                    = "customer.created"
  CustomerUpdatedEvent                    = "customer.updated"
  CustomerStatusUpdatedEvent              = "customer.status.updated"
  CustomerBalanceUpdatedEvent             = "customer.balance.updated"
  CustomerPasswordUpdatedEvent            = "customer.password.updated"
  CustomerPhotoUpdatedEvent               = "customer.avatar.updated"
  CustomerEmailVerifiedEvent              = "customer.email.verified"
  CustomerForgotPasswordRequestedEvent    = "customer.forgot-password.requested"
  CustomerEmailVerificationRequestedEvent = "customer.email-verif.requested"
  CustomerAddressAddedEvent               = "customer.address.added"
  CustomerAddressDeletedEvent             = "customer.addresses.deleted"
  CustomerAddressUpdatedEvent             = "customer.address.updated"
  CustomerVoucherAddedEvent               = "customer.vouchers.added"
  CustomerVoucherDeletedEvent             = "customer.vouchers.deleted"
  CustomerVoucherUpdatedEvent             = "customer.voucher.updated"
  CustomerDefaultAddressUpdatedEvent      = "customer.default-address.updated"
)
