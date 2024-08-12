package types

import (
  "github.com/google/uuid"
  "strconv"
  "time"
)

type Keyer interface {
  // Key return if the event has key, it will return false for bool if it has no key
  Key() (string, bool)
}

type Metadataer interface {
  Metadata() Metadata
}

type Event interface {
  Entity
  Keyer
  //Metadataer
  EventVersioning
  EventTyping
  EventName() string
  OccurredAt() time.Time
}

type EventOption[T EventTyping, V EventVersioning] func(base *EventBase[T, V])

// WithTime set custom time for occurredAt field instead of time.Now() as defaulted one
func WithTime[T EventTyping, V EventVersioning](t time.Time) EventOption[T, V] {
  return func(base *EventBase[T, V]) {
    base.occurredAt = t
  }
}

// WithId set the id instead of using the default which is using uuid
func WithId[T EventTyping, V EventVersioning](id string) EventOption[T, V] {
  return func(base *EventBase[T, V]) {
    base.id = id
  }
}

func NewEvent[T EventTyping, V EventVersioning](opts ...EventOption[T, V]) EventBase[T, V] {
  ev := EventBase[T, V]{
    occurredAt: time.Now(),
    id:         uuid.NewString(),
  }

  for _, opt := range opts {
    opt(&ev)
  }
  return ev
}

// EventBase base of event without typing, versioning and empty key.
// Override the method interface to provide them or use EventBaseV1 to add v1 as versioning,
// DomainEventBase to add EventTypeDomain as event type and event DomainEventBaseV1 to add default
// EventTypeDomain as event type and 1 for the version
type EventBase[T EventTyping, V EventVersioning] struct {
  types      T
  version    V
  occurredAt time.Time
  id         string // Unique per event
}

func (e *EventBase[T, V]) OccurredAt() time.Time {
  return e.occurredAt
}

func (e *EventBase[T, V]) Identity() string {
  return e.id
}

func (e *EventBase[T, V]) Key() (string, bool) {
  return "", false
}

func (e *EventBase[T, V]) EventVersion() uint8 {
  return e.version.EventVersion()
}

func (e *EventBase[T, V]) EventType() EventType {
  return e.types.EventType()
}

func ConstructMetadata(e Event) Metadata {
  return Metadata{
    NewKeyVal(METADATA_IDENTITY_KEY, e.Identity()),
    NewKeyVal(METADATA_EVENT_NAME_KEY, e.EventName()),
    NewKeyVal(METADATA_VERSION_KEY, strconv.Itoa(int(e.EventVersion()))),
    NewKeyVal(METADATA_EVENT_TYPE_KEY, string(e.EventType())),
  }
}

type DomainEvent = EventBase[DomainEventType, NoVersioning]

func NewDomainEvent[V EventVersioning](options ...EventOption[DomainEventType, V]) EventBase[DomainEventType, V] {
  return NewEvent[DomainEventType, V](options...)
}

type DomainEventV1 = EventBase[DomainEventType, V1]

type IntegrationEvent = EventBase[IntegrationEventType, NoVersioning]

func NewIntegrationEvent[V EventVersioning](options ...EventOption[IntegrationEventType, V]) EventBase[IntegrationEventType, V] {
  return NewEvent[IntegrationEventType, V](options...)
}

type IntegrationEventV1 = EventBase[IntegrationEventType, V1]
