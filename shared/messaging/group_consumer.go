package messaging

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/logger"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
  "github.com/dnwe/otelsarama"
)

func NewGroupConsumerBase[T any](handler T, deserializer serde.IDeserializerAny, forwarder IConsumerForwarder) GroupConsumerBase[T] {
  return GroupConsumerBase[T]{
    handler:      handler,
    deserializer: deserializer,
    errors:       make(chan errorMessage),
    dlqPublisher: forwarder,
  }
}

var _ sarama.ConsumerGroupHandler = (*GroupConsumerBase[any])(nil)

type errorMessage struct {
  id      string
  event   string
  message *sarama.ConsumerMessage
  error   error
}

type GroupConsumerBase[T any] struct {
  handler      T
  deserializer serde.IDeserializerAny

  errors       chan errorMessage
  dlqPublisher IConsumerForwarder
}

func (g *GroupConsumerBase[T]) Handler() T {
  return g.handler
}

func (g *GroupConsumerBase[T]) Deserializer() serde.IDeserializerAny {
  return g.deserializer
}

func (g *GroupConsumerBase[T]) Run(ctx context.Context, group sarama.ConsumerGroup, topics ...string) {
  go func() {
    groupHandler := otelsarama.WrapConsumerGroupHandler(g)
    for {
      err := group.Consume(ctx, topics, groupHandler)

      if err != nil {
        if ctx.Err() != nil {
          return
        }
        logger.Infof("Failed to consume message from topic %s: %v", topics, err)
      }
      logger.Infof("Successfully consumed message from topic %s", topics)
    }
  }()
}

// errorHandler is default handler for each error on consuming or processing message or event.
// It will publish or forward the message into DLQ topic
func (g *GroupConsumerBase[T]) errorHandler() {
  for msg := range g.errors {
    logger.Infof("[%s] Error on processing %s::%s message, sent to DLQ", string(msg.message.Key), msg.event, msg.id)

    // Send to DLQ
    if g.dlqPublisher == nil {
      continue
    }
    err := g.dlqPublisher.Forward(context.Background(), msg.message, msg.error)
    if err != nil {
      logger.Infof("Failed to forward message to DLQ: %v", err)
    }
  }
}

func (g *GroupConsumerBase[T]) Setup(session sarama.ConsumerGroupSession) error {
  go g.errorHandler()
  return nil
}

func (g *GroupConsumerBase[T]) handleError(metadata *types.EventMetadata, message *sarama.ConsumerMessage, err error) {
  g.errors <- errorMessage{
    id:      metadata.Id,
    event:   metadata.Name,
    message: message,
    error:   err,
  }
}

func (g *GroupConsumerBase[T]) Cleanup(session sarama.ConsumerGroupSession) error {
  close(g.errors)
  return nil
}

// ProcessMessage should be implemented by the child, or it will be noop
func (g *GroupConsumerBase[T]) ProcessMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage, md *types.EventMetadata) error {
  return nil
}

func (g *GroupConsumerBase[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
  defer func() {
    if r := recover(); r != nil {
      logger.Infof("GroupConsumerBase::ConsumeClaim - Recovered from panic: %v", r)
    }
  }()

  logger.Infof("Start consuming %s on [%d] starting at %d offset", claim.Topic(), claim.Partition(), claim.InitialOffset())

  for message := range claim.Messages() {
    logger.Infof("Got Message %s::%s ", message.Topic, message.Key)

    md := types.NewEventMetadata(message.Headers, types.FromKafkaRecords())
    err := g.ProcessMessage(session, message, md)
    if err != nil {
      g.handleError(md, message, err)
    }
  }
  return nil
}
