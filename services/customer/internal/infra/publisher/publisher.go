package publisher

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "mini-shop/services/user/internal/domain/repository"
  "mini-shop/services/user/pkg/tracer"
)

func NewKafka(producer sarama.SyncProducer, serializer serde.ISerializer) repository.IMessagePublisher {
  return &kafkaPublisher{
    producer:   producer,
    serializer: serializer,
    tracer:     tracer.Get(),
  }
}

type kafkaPublisher struct {
  producer   sarama.SyncProducer
  serializer serde.ISerializer
  tracer     trace.Tracer
}

func (k *kafkaPublisher) Close() error {
  return k.producer.Close()
}

func (k *kafkaPublisher) GracefulShutdown(_ context.Context) error {
  return k.Close()
}

func (k *kafkaPublisher) PublishEvents(ctx context.Context, events ...types.Event) error {
  ctx, span := k.tracer.Start(ctx, "kafkaPublisher.PublishEvents")
  defer span.End()

  if len(events) == 0 {
    return nil
  }

  messages, ierr := sharedUtil.CastSliceErrs(events, func(event types.Event) (*sarama.ProducerMessage, error) {
    bytes, err := k.serializer.Serialize(event)
    if err != nil {
      return nil, err
    }

    return &sarama.ProducerMessage{
      Topic:     event.EventName(),
      Key:       sarama.StringEncoder(event.Key()),
      Value:     sarama.ByteEncoder(bytes),
      Headers:   event.Metadata().ToKafkaRecordHeader(),
      Timestamp: event.OccurredAt(),
    }, nil
  })

  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return ierr
  }

  err := k.producer.SendMessages(messages)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}

func (k *kafkaPublisher) PublishAggregate(ctx context.Context, aggregate types.Aggregate) error {
  ctx, span := k.tracer.Start(ctx, "kafkaPublisher.PublishAggregate")
  defer span.End()

  return k.PublishEvents(ctx, aggregate.Events()...)
}
