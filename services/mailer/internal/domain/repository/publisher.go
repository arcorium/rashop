package repository

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
  "github.com/arcorium/rashop/shared/types"
)

type IMessagePublisher interface {
  PublishEvents(ctx context.Context, events ...types.Event) error
  Publish(ctx context.Context, customer *entity.Mail) error
  // Close will cut all not-delivered messages
  Close() error
}
