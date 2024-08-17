package messaging

import (
  "context"
  "github.com/arcorium/rashop/shared/types"
)

type IPublisher interface {
  PublishEvents(ctx context.Context, events ...types.Event) error
}
