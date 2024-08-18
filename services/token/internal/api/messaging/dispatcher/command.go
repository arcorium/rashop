package dispatcher

import (
  "fmt"
  "github.com/IBM/sarama"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/token/internal/api/messaging/consumer"
  "github.com/arcorium/rashop/shared/messaging"
  "github.com/arcorium/rashop/shared/otel"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
)

func NewTokenCommandConsumerGroup(dlqPublisher messaging.IForwarder[*sarama.ConsumerMessage], handler consumer.ITokenCommandHandler, deserializer serde.IDeserializerAny) *TokenCommandConsumerGroupHandler {
  return &TokenCommandConsumerGroupHandler{
    GroupConsumerBase: messaging.NewGroupConsumerBase(handler, deserializer, dlqPublisher),
  }
}

type TokenCommandConsumerGroupHandler struct {
  messaging.GroupConsumerBase[consumer.ITokenCommandHandler]
}

func (q *TokenCommandConsumerGroupHandler) ProcessMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage, md *types.EventMetadata) error {
  var err error
  // Parse OTEL from message
  ctx := otel.Extract(md, message)

  switch md.Name {
  case intev.CustomerCreatedEvent:
    err = messaging.DispatchV1[*intev.CustomerCreatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&intev.CustomerCreatedV1{}, md.Id, message.Timestamp), q.Handler().OnCustomerCreated)
  case intev.CustomerEmailVerificationRequestedEvent:
    err = messaging.DispatchV1[*intev.CustomerEmailVerificationRequestedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&intev.CustomerEmailVerificationRequestedV1{}, md.Id, message.Timestamp), q.Handler().OnEmailVerificationRequested)
  case intev.LoginTokenRequestedEvent:
    err = messaging.DispatchV1[*intev.LoginTokenRequestedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&intev.LoginTokenRequestedV1{}, md.Id, message.Timestamp), q.Handler().OnLoginRequested)
  case intev.CustomerResetPasswordRequestedEvent:
    err = messaging.DispatchV1[*intev.CustomerResetPasswordRequestedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&intev.CustomerResetPasswordRequestedV1{}, md.Id, message.Timestamp), q.Handler().OnResetPasswordRequested)
  case intev.CustomerEmailUpdatedEvent:
    err = messaging.DispatchV1[*intev.CustomerEmailUpdatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&intev.CustomerEmailUpdatedV1{}, md.Id, message.Timestamp), q.Handler().OnCustomerEmailUpdated)
  default:
    err = fmt.Errorf("unknown event : %s", md.Name)
  }

  session.MarkMessage(message, "")
  return err
}
