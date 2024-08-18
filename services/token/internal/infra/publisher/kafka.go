package publisher

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/services/token/internal/domain/entity"
  "github.com/arcorium/rashop/services/token/internal/domain/repository"
  "github.com/arcorium/rashop/services/token/pkg/tracer"
  "github.com/arcorium/rashop/shared/messaging/kafka"
  "github.com/arcorium/rashop/shared/serde"
)

func NewKafka(topic KafkaTopic, producer sarama.SyncProducer, serializer serde.ISerializer) repository.IMessagePublisher {
  return &tokenKafkaPublisher{
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

type tokenKafkaPublisher struct {
  kafka.PublisherBase
}

func (m *tokenKafkaPublisher) Publish(ctx context.Context, token *entity.Token) error {
  ctx, span := m.Tracer.Start(ctx, "tokenKafkaPublisher.Publish")
  defer span.End()

  return m.PublishEvents(ctx, token.Events()...)
}
