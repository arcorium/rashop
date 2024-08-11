package handler

import (
  "context"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
)

type Command[C any, R any] interface {
  Handle(context.Context, C) (R, status.Object)
}

type CommandUnit[C any] interface {
  Handle(context.Context, C) status.Object
}

type Consumer[E any] interface {
  Handle(context.Context, E) status.Object
}

type ConsumerFunc[E types.Event] func(context.Context, E) error
