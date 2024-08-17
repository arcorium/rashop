package publisher

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
  "github.com/arcorium/rashop/services/mailer/internal/domain/repository"
  "github.com/arcorium/rashop/services/mailer/pkg/tracer"
  "github.com/arcorium/rashop/shared/messaging/kafka"
  "github.com/arcorium/rashop/shared/serde"
)

func NewKafka(topic KafkaTopic, producer sarama.SyncProducer, serializer serde.ISerializer) repository.IMessagePublisher {
  return &mailerKafkaPublisher{
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

type mailerKafkaPublisher struct {
  kafka.PublisherBase
}

func (k *mailerKafkaPublisher) Publish(ctx context.Context, customer *entity.Mail) error {
  ctx, span := k.Tracer.Start(ctx, "mailerKafkaPublisher.Publish")
  defer span.End()

  return k.PublishEvents(ctx, customer.Events()...)
}
