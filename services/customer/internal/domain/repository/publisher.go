package repository

import (
  "context"
  "github.com/arcorium/rashop/shared/types"
  "rashop/services/customer/internal/domain/entity"
)

type IMessagePublisher interface {
  PublishEvents(ctx context.Context, events ...types.Event) error
  Publish(ctx context.Context, customer *entity.Customer) error
  // Close will cut all not-delivered messages
  Close() error
}
