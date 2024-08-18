package repository

import (
  "context"
  "github.com/arcorium/rashop/services/token/internal/domain/entity"
  "github.com/arcorium/rashop/shared/types"
)

type IMessagePublisher interface {
  PublishEvents(ctx context.Context, events ...types.Event) error
  Publish(ctx context.Context, token *entity.Token) error
}
