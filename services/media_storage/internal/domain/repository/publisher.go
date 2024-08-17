package repository

import (
  "context"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/entity"
  "github.com/arcorium/rashop/shared/types"
)

type IMessagePublisher interface {
  PublishEvents(ctx context.Context, events ...types.Event) error
  Publish(ctx context.Context, media *entity.Media) error
}
