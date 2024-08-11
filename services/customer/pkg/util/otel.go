package util

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/types"
  "github.com/dnwe/otelsarama"
  "go.opentelemetry.io/otel"
)

func GetOTELFromKafkaMessage(parent context.Context, message *sarama.ConsumerMessage) context.Context {
  return otel.GetTextMapPropagator().
    Extract(parent, otelsarama.NewConsumerMessageCarrier(message))
}

func ExtractOTEL(metadata *types.EventMetadata, msg *sarama.ConsumerMessage) context.Context {
  return GetOTELFromKafkaMessage(metadata.AsContext(context.Background()), msg)
}

func InjectOTEL(ctx context.Context, message *sarama.ProducerMessage) {
  otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(message))
}
