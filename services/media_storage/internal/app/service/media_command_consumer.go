package service

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/media_storage/internal/app/command/consumer"
  "github.com/arcorium/rashop/shared/status"
)

type IMediaCommandConsumer interface {
  OneTimeMediaUsed(ctx context.Context, v1 *intev.OneTimeMediaUsedV1) status.Object
}

func NewMediaCommandConsumer(config MediaCommandConsumerConfig) IMediaCommandConsumer {
  return &mediaConsumerService{
    i: config,
  }
}

func DefaultMediaCommandConsumerConfig(parameter consumer.CommonHandlerParameter) MediaCommandConsumerConfig {
  return MediaCommandConsumerConfig{
    OneTimeUsed: consumer.NewOneTimeMediaUsedHandler(parameter),
  }
}

type MediaCommandConsumerConfig struct {
  OneTimeUsed consumer.IOneTimeMediaUsedHandler
}

type mediaConsumerService struct {
  i MediaCommandConsumerConfig
}

func (m *mediaConsumerService) OneTimeMediaUsed(ctx context.Context, ev *intev.OneTimeMediaUsedV1) status.Object {
  return m.i.OneTimeUsed.Handle(ctx, ev)
}
