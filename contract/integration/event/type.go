package intev

import "github.com/arcorium/rashop/shared/types"

type IntegrationV1 = types.EventBase[types.IntegrationEventType, types.V1]

func NewV1(options ...types.EventOption[types.IntegrationEventType, types.V1]) IntegrationV1 {
  return types.NewIntegrationEvent[types.V1](options...)
}
