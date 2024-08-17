package dispatcher

import (
  "fmt"
  "github.com/IBM/sarama"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/media_storage/internal/api/messaging/consumer"
  "github.com/arcorium/rashop/shared/messaging"
  "github.com/arcorium/rashop/shared/otel"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
)

func NewMediaCommandConsumerGroup(dlqPublisher messaging.IConsumerForwarder, handler consumer.IMediaCommandHandler, deser serde.IDeserializerAny) *MediaCommandConsumerGroupHandler {
  return &MediaCommandConsumerGroupHandler{
    GroupConsumerBase: messaging.NewGroupConsumerBase(handler, deser, dlqPublisher),
  }
}

type MediaCommandConsumerGroupHandler struct {
  messaging.GroupConsumerBase[consumer.IMediaCommandHandler]
}

func (q *MediaCommandConsumerGroupHandler) ProcessMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage, md *types.EventMetadata) error {
  var err error
  // Parse OTEL from message
  ctx := otel.Extract(md, message)

  switch md.Name {
  case intev.MediaStorageOneTimeUsedEvent:
    err = messaging.DispatchV1[*intev.OneTimeMediaUsedV1](q.Deserializer(), ctx, message,
      types.ConstructEventBase(&intev.OneTimeMediaUsedV1{}, md.Id, message.Timestamp),
      q.Handler().OnOneTimeMediaUsed)
  default:
    err = fmt.Errorf("unknown event : %s", md.Name)
  }

  session.MarkMessage(message, "")
  return err
}
