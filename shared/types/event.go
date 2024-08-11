package types

import (
  "github.com/google/uuid"
  "strconv"
  "time"
)

type Keyer interface {
  Key() string
}

type Metadataer interface {
  Metadata() Metadata
}

type Event interface {
  Entity
  Keyer
  Metadataer
  EventVersioning
  EventTyping
  EventName() string
  OccurredAt() time.Time
}

type EventOption func(base *EventBase)

// WithTime set custom time for occurredAt field instead of time.Now() as defaulted one
func WithTime(t time.Time) EventOption {
  return func(base *EventBase) {
    base.occurredAt = t
  }
}

// WithId set the id instead of using the default which is using uuid
func WithId(id string) EventOption {
  return func(base *EventBase) {
    base.id = id
  }
}

func NewEvent(opts ...EventOption) EventBase {
  event := EventBase{
    id:         uuid.NewString(),
    occurredAt: time.Now(),
  }

  for _, opt := range opts {
    opt(&event)
  }
  return event
}

// EventBase base of event without typing, versioning and empty key.
// Override the method interface to provide them or use EventBaseV1 to add v1 as versioning,
// DomainEventBase to add EventTypeDomain as event type and event DomainEventBaseV1 to add default
// EventTypeDomain as event type and 1 for the version
type EventBase struct {
  NoVersioning
  NoEventTyping
  occurredAt time.Time
  id         string // Unique per event
}

func (e *EventBase) OccurredAt() time.Time {
  return e.occurredAt
}

func (e *EventBase) Identity() string {
  return e.id
}

func (e *EventBase) Key() string {
  return ""
}

func (e *EventBase) Metadata() Metadata {
  return Metadata{
    NewKeyVal(METADATA_IDENTITY_KEY, e.id),
    NewKeyVal(METADATA_VERSION_KEY, strconv.Itoa(int(e.EventVersion()))),
    NewKeyVal(METADATA_EVENT_TYPE_KEY, string(e.EvenType())),
  }
}

func NewDomainEvent(options ...EventOption) DomainEventBase {
  return DomainEventBase{
    EventBase: NewEvent(options...),
  }
}

type DomainEventBase struct {
  EventBase
}

func (d *DomainEventBase) EventType() EventType {
  return EventTypeDomain
}

func NewIntegrationEvent(options ...EventOption) IntegrationEventBase {
  return IntegrationEventBase{
    EventBase: NewEvent(options...),
  }
}

type IntegrationEventBase struct {
  EventBase
}

func (d *IntegrationEventBase) EventType() EventType {
  return EventTypeIntegration
}

func NewEventV1(options ...EventOption) EventBaseV1 {
  return EventBaseV1{
    EventBase: NewEvent(options...),
    V1:        V1{},
  }
}

type EventBaseV1 struct {
  EventBase
  V1
}

func NewEventV2(options ...EventOption) EventBaseV2 {
  return EventBaseV2{
    EventBase: NewEvent(options...),
    V2:        V2{},
  }
}

type EventBaseV2 struct {
  EventBase
  V2
}

func NewDomainEventV1(options ...EventOption) DomainEventBaseV1 {
  return DomainEventBaseV1{
    DomainEventBase: NewDomainEvent(options...),
    V1:              V1{},
  }
}

type Constructor[T Event] func(*T)

type DomainEventBaseV1 struct {
  DomainEventBase
  V1
}

func NewIntegrationEventV1(options ...EventOption) IntegrationEventBaseV1 {
  return IntegrationEventBaseV1{
    IntegrationEventBase: NewIntegrationEvent(options...),
    V1:                   V1{},
  }
}

type IntegrationEventBaseV1 struct {
  IntegrationEventBase
  V1
}
