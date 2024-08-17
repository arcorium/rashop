package kafka

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
)

func NewPublisherBase(producer sarama.SyncProducer, serializer serde.ISerializer, tracer trace.Tracer, domainTopic, integrateTopic string) PublisherBase {
  return PublisherBase{
    producer:              producer,
    serializer:            serializer,
    Tracer:                tracer,
    domainEventTopic:      domainTopic,
    integrationEventTopic: integrateTopic,
  }
}

type PublisherBase struct {
  producer   sarama.SyncProducer
  serializer serde.ISerializer
  Tracer     trace.Tracer

  domainEventTopic      string
  integrationEventTopic string
}

func (p *PublisherBase) PublishEvents(ctx context.Context, events ...types.Event) error {
  ctx, span := p.Tracer.Start(ctx, "PublisherBase.PublishEvents")
  defer span.End()

  if len(events) == 0 {
    return nil
  }

  messages, ierr := sharedUtil.CastSliceErrs(events, func(event types.Event) (*sarama.ProducerMessage, error) {
    // Serialize value
    bytes, err := p.serializer.Serialize(event)
    if err != nil {
      return nil, err
    }

    // GetCustomers key
    var key sarama.Encoder
    keys, ok := event.Key()
    if ok {
      key = sarama.StringEncoder(keys)
    }

    // Construct metadata as header
    metadata := types.ConstructMetadata(event)

    // Determine topic
    topic := p.domainEventTopic
    if event.EventType() == types.EventTypeIntegration {
      topic = p.integrationEventTopic
    }

    return &sarama.ProducerMessage{
      Topic:     topic,
      Key:       key,
      Value:     sarama.ByteEncoder(bytes),
      Headers:   metadata.ToKafkaRecordHeader(),
      Timestamp: event.OccurredAt(),
    }, nil
  })

  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return ierr
  }

  err := p.producer.SendMessages(messages)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}
