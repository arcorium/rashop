package event

import "github.com/arcorium/rashop/shared/types"

type DomainV1 = types.EventBase[types.DomainEventType, types.V1]

func NewV1(options ...types.EventOption[types.DomainEventType, types.V1]) DomainV1 {
  return types.NewDomainEvent[types.V1](options...)
}
