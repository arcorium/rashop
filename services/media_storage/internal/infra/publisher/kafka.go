package publisher

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/entity"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/repository"
  "github.com/arcorium/rashop/services/media_storage/pkg/tracer"
  "github.com/arcorium/rashop/shared/messaging/kafka"
  "github.com/arcorium/rashop/shared/serde"
)

func NewKafka(topic KafkaTopic, producer sarama.SyncProducer, serializer serde.ISerializer) repository.IMessagePublisher {
  return &mediaKafkaPublisher{
    PublisherBase: kafka.NewPublisherBase(producer,
      serializer,
      tracer.Get(),
      topic.DomainEvent,
      topic.IntegrationEvent),
  }
}

type KafkaTopic struct {
  DomainEvent      string
  IntegrationEvent string
}

type mediaKafkaPublisher struct {
  kafka.PublisherBase
}

func (k *mediaKafkaPublisher) Publish(ctx context.Context, customer *entity.Media) error {
  ctx, span := k.Tracer.Start(ctx, "mediaKafkaPublisher.Publish")
  defer span.End()

  return k.PublishEvents(ctx, customer.Events()...)
}
