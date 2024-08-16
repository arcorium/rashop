package dispatcher

import (
  "fmt"
  "github.com/IBM/sarama"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/mailer/internal/api/messaging/consumer"
  "github.com/arcorium/rashop/shared/messaging"
  "github.com/arcorium/rashop/shared/otel"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
)

func NewMailCommandConsumerGroup(dlqPublisher messaging.IForwarder[*sarama.ConsumerMessage], handler consumer.IMailCommandHandler, deserializer serde.IDeserializerAny) *MailerCommandConsumerGroupHandler {
  return &MailerCommandConsumerGroupHandler{
    GroupConsumerBase: messaging.NewGroupConsumerBase(handler, deserializer, dlqPublisher),
  }
}

type MailerCommandConsumerGroupHandler struct {
  messaging.GroupConsumerBase[consumer.IMailCommandHandler]
}

func (q *MailerCommandConsumerGroupHandler) ProcessMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage, md *types.EventMetadata) error {
  var err error
  // Parse OTEL from message
  ctx := otel.Extract(md, message)

  switch md.Name {
  case intev.EmailVerificationTokenCreatedEvent:
    err = messaging.DispatchV1[*intev.EmailVerificationTokenCreatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&intev.EmailVerificationTokenCreatedV1{}, md.Id, message.Timestamp), q.Handler().OnEmailVerificationTokenCreated)
  case intev.ResetPasswordTokenCreatedEvent:
    err = messaging.DispatchV1[*intev.ResetPasswordTokenCreatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&intev.ResetPasswordTokenCreatedV1{}, md.Id, message.Timestamp), q.Handler().OnResetPasswordTokenCreated)
  case intev.LoginTokenCreatedEvent:
    err = messaging.DispatchV1[*intev.LoginTokenCreatedV1](q.Deserializer(), ctx, message, types.ConstructEventBase(&intev.LoginTokenCreatedV1{}, md.Id, message.Timestamp), q.Handler().OnLoginTokenCreated)
  default:
    err = fmt.Errorf("unknown event : %s", md.Name)
  }

  session.MarkMessage(message, "")
  return err
}
