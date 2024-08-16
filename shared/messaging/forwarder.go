package messaging

import (
  "context"
  "github.com/IBM/sarama"
)

const HEADER_ERROR_KEY = "error"

type IForwarder[T any] interface {
  Forward(ctx context.Context, message T, err error) error
}

type IConsumerForwarder = IForwarder[*sarama.ConsumerMessage]
