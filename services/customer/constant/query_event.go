package constant

import "mini-shop/services/user/internal/domain/event"

// TODO: Add integration event
var CustomerQuerySubscribedDomainEvent = []string{
  event.CustomerCreatedEvent,
  event.CustomerUpdatedEvent,
  event.CustomerStatusUpdatedEvent,
  event.CustomerBalanceUpdatedEvent,
  event.CustomerPasswordUpdatedEvent,
  event.CustomerPhotoUpdatedEvent,
  event.CustomerEmailVerifiedEvent,
  event.CustomerDeletedEvent,
  event.CustomerForgotPasswordRequestedEvent,
  event.CustomerEmailVerificationRequestedEvent,
  event.CustomerAddressAddedEvent,
  event.CustomerAddressesDeletedEvent,
  event.CustomerAddressUpdatedEvent,
  event.CustomerVouchersAddedEvent,
  event.CustomerVouchersDeletedEvent,
  event.CustomerVoucherUpdatedEvent,
  event.CustomerDefaultAddressUpdatedEvent,
}
