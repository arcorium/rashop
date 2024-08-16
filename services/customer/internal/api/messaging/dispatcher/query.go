package dispatcher

import (
  "fmt"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/messaging"
  "github.com/arcorium/rashop/shared/otel"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
  "rashop/services/customer/internal/api/messaging/consumer"
  "rashop/services/customer/internal/domain/event"
)

func NewQueryConsumerGroup(dlqPublisher messaging.IForwarder[*sarama.ConsumerMessage], handler consumer.ICustomerQueryHandler, deserializer serde.IDeserializerAny) *QueryConsumerGroupHandler {
  return &QueryConsumerGroupHandler{
    GroupConsumerBase: messaging.NewGroupConsumerBase(handler, deserializer, dlqPublisher),
  }
}

type QueryConsumerGroupHandler struct {
  messaging.GroupConsumerBase[consumer.ICustomerQueryHandler]
}

func (q *QueryConsumerGroupHandler) ProcessMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage, md *types.EventMetadata) error {
  var err error
  // Parse OTEL from message
  ctx := otel.Extract(md, message)

  switch md.Name {
  case event.CustomerCreatedEvent:
    err = messaging.DispatchV1[*event.CustomerCreatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerCreatedV1{}, md.Id, message.Timestamp), q.Handler().OnCreatedV1)
  case event.CustomerUpdatedEvent:
    err = messaging.DispatchV1[*event.CustomerUpdatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerUpdatedV1{}, md.Id, message.Timestamp), q.Handler().OnUpdatedV1)
  case event.CustomerStatusUpdatedEvent:
    err = messaging.DispatchV1[*event.CustomerStatusUpdatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerStatusUpdatedV1{}, md.Id, message.Timestamp), q.Handler().OnStatusUpdatedV1)
  case event.CustomerBalanceUpdatedEvent:
    err = messaging.DispatchV1[*event.CustomerBalanceUpdatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerBalanceUpdatedV1{}, md.Id, message.Timestamp), q.Handler().OnBalanceUpdatedV1)
  case event.CustomerPasswordUpdatedEvent:
    err = messaging.DispatchV1[*event.CustomerPasswordUpdatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerPasswordUpdatedV1{}, md.Id, message.Timestamp), q.Handler().OnPasswordUpdatedV1)
  case event.CustomerPhotoUpdatedEvent:
    err = messaging.DispatchV1[*event.CustomerPhotoUpdatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerPhotoUpdatedV1{}, md.Id, message.Timestamp), q.Handler().OnPhotoUpdatedV1)
  case event.CustomerEmailVerifiedEvent:
    err = messaging.DispatchV1[*event.CustomerEmailVerifiedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerEmailVerifiedV1{}, md.Id, message.Timestamp), q.Handler().OnEmailVerifiedV1)
  case event.CustomerAddressAddedEvent:
    err = messaging.DispatchV1[*event.CustomerAddressAddedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerAddressAddedV1{}, md.Id, message.Timestamp), q.Handler().OnAddressAddedV1)
  case event.CustomerAddressDeletedEvent:
    err = messaging.DispatchV1[*event.CustomerAddressDeletedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerAddressDeletedV1{}, md.Id, message.Timestamp), q.Handler().OnAddressDeletedV1)
  case event.CustomerAddressUpdatedEvent:
    err = messaging.DispatchV1[*event.CustomerAddressUpdatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerAddressUpdatedV1{}, md.Id, message.Timestamp), q.Handler().OnAddressUpdatedV1)
  case event.CustomerVoucherAddedEvent:
    err = messaging.DispatchV1[*event.CustomerVoucherAddedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerVoucherAddedV1{}, md.Id, message.Timestamp), q.Handler().OnVoucherAddedV1)
  case event.CustomerVoucherDeletedEvent:
    err = messaging.DispatchV1[*event.CustomerVoucherDeletedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerVoucherDeletedV1{}, md.Id, message.Timestamp), q.Handler().OnVoucherDeletedV1)
  case event.CustomerVoucherUpdatedEvent:
    err = messaging.DispatchV1[*event.CustomerVoucherUpdatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerVoucherUpdatedV1{}, md.Id, message.Timestamp), q.Handler().OnVoucherUpdatedV1)
  case event.CustomerDefaultAddressUpdatedEvent:
    err = messaging.DispatchV1[*event.CustomerDefaultAddressUpdatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.CustomerDefaultAddressUpdatedV1{}, md.Id, message.Timestamp), q.Handler().OnDefaultAddressUpdatedV1)
  default:
    err = fmt.Errorf("unknown event : %s", md.Name)
  }

  session.MarkMessage(message, "") // Mark as processed either there is an error or not
  return err
}
