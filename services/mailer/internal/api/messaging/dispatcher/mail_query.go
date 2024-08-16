package dispatcher

import (
  "fmt"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/services/mailer/internal/api/messaging/consumer"
  "github.com/arcorium/rashop/services/mailer/internal/domain/event"
  "github.com/arcorium/rashop/shared/messaging"
  "github.com/arcorium/rashop/shared/otel"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
)

func NewMailQueryConsumerGroup(dlqPublisher messaging.IForwarder[*sarama.ConsumerMessage], handler consumer.IMailQueryHandler, deserializer serde.IDeserializerAny) *MailerQueryConsumerGroupHandler {
  return &MailerQueryConsumerGroupHandler{
    GroupConsumerBase: messaging.NewGroupConsumerBase(handler, deserializer, dlqPublisher),
  }
}

type MailerQueryConsumerGroupHandler struct {
  messaging.GroupConsumerBase[consumer.IMailQueryHandler]
}

func (q *MailerQueryConsumerGroupHandler) ProcessMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage, md *types.EventMetadata) error {
  var err error
  // Parse OTEL from message
  ctx := otel.Extract(md, message)

  switch md.Name {
  case event.MailerCreatedEvent:
    err = messaging.DispatchV1[*event.MailCreatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.MailCreatedV1{}, md.Id, message.Timestamp), q.Handler().OnMailCreated)
  case event.MailerStatusUpdatedEvent:
    err = messaging.DispatchV1[*event.MailStatusUpdatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.MailStatusUpdatedV1{}, md.Id, message.Timestamp), q.Handler().OnMailStatusUpdated)
  case event.MailerDeletedEvent:
    err = messaging.DispatchV1[*event.MailDeletedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&event.MailDeletedV1{}, md.Id, message.Timestamp), q.Handler().OnMailDeleted)
  default:
    err = fmt.Errorf("unknown event : %s", md.Name)
  }

  session.MarkMessage(message, "")
  return err
}
