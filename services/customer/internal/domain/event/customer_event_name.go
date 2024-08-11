package event

const (
  CustomerCreatedEvent                    = "dom.customer.created"
  CustomerUpdatedEvent                    = "dom.customer.updated"
  CustomerStatusUpdatedEvent              = "dom.customer.status.updated"
  CustomerBalanceUpdatedEvent             = "dom.customer.balance.updated"
  CustomerPasswordUpdatedEvent            = "dom.customer.password.updated"
  CustomerPhotoUpdatedEvent               = "dom.customer.avatar.updated"
  CustomerEmailVerifiedEvent              = "dom.customer.email.verified"
  CustomerDeletedEvent                    = "dom.customer.disabled"
  CustomerForgotPasswordRequestedEvent    = "dom.customer.forgot-password.requested"
  CustomerEmailVerificationRequestedEvent = "dom.customer.email-verif.requested"
  CustomerAddressAddedEvent               = "dom.customer.address.added"
  CustomerAddressesDeletedEvent           = "dom.customer.addresses.deleted"
  CustomerAddressUpdatedEvent             = "dom.customer.address.updated"
  CustomerVouchersAddedEvent              = "dom.customer.vouchers.added"
  CustomerVouchersDeletedEvent            = "dom.customer.vouchers.deleted"
  CustomerVoucherUpdatedEvent             = "dom.customer.voucher.updated"
  CustomerDefaultAddressUpdatedEvent      = "dom.customer.default-address.updated"
)
