package repository

import (
  "context"
  "github.com/arcorium/rashop/shared/types"
)

type IMessagePublisher interface {
  PublishEvents(ctx context.Context, events ...types.Event) error
  PublishAggregate(ctx context.Context, aggregate types.Aggregate) error
  // Close will cut all not-delivered messages
  Close() error
  // GracefulShutdown will waits until all messages delivered and close it afterward
  GracefulShutdown(ctx context.Context) error
}
