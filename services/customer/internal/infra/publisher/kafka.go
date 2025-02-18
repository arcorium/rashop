package publisher

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/messaging/kafka"
  "github.com/arcorium/rashop/shared/serde"
  "rashop/services/customer/internal/domain/entity"
  "rashop/services/customer/internal/domain/repository"
  "rashop/services/customer/pkg/tracer"
)

func NewKafka(topic KafkaTopic, producer sarama.SyncProducer, serializer serde.ISerializer) repository.IMessagePublisher {
  return &customerKafkaPublisher{
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

type customerKafkaPublisher struct {
  kafka.PublisherBase
}

func (k *customerKafkaPublisher) Publish(ctx context.Context, customer *entity.Customer) error {
  ctx, span := k.Tracer.Start(ctx, "customerKafkaPublisher.PublishAggregate")
  defer span.End()

  return k.PublishEvents(ctx, customer.Events()...)
}
