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

type EventOption[T EventTyping, V EventVersioning] func(base *eventBase[T, V])

// WithTime set custom time for occurredAt field instead of time.Now() as defaulted one
func WithTime[T EventTyping, V EventVersioning](t time.Time) EventOption[T, V] {
  return func(base *eventBase[T, V]) {
    base.occurredAt = t
  }
}

// WithId set the id instead of using the default which is using uuid
func WithId[T EventTyping, V EventVersioning](id string) EventOption[T, V] {
  return func(base *eventBase[T, V]) {
    base.id = id
  }
}

func newEvent[T EventTyping, V EventVersioning](opts ...EventOption[T, V]) eventBase[T, V] {
  ev := eventBase[T, V]{
    occurredAt: time.Now(),
    id:         uuid.NewString(),
  }

  for _, opt := range opts {
    opt(&ev)
  }
  return ev
}

type IEventBaseConstructor interface {
}

// eventBase base of event without typing, versioning and empty key.
// Override the method interface to provide them or use EventBaseV1 to add v1 as versioning,
// DomainEventBase to add EventTypeDomain as event type and event DomainEventBaseV1 to add default
// EventTypeDomain as event type and 1 for the version
type eventBase[T EventTyping, V EventVersioning] struct {
  types      T
  version    V
  occurredAt time.Time
  id         string // Unique per event
}

func (e *eventBase[T, V]) OccurredAt() time.Time {
  return e.occurredAt
}

func (e *eventBase[T, V]) Identity() string {
  return e.id
}

func (e *eventBase[T, V]) Key() (string, bool) {
  return "", false
}

func (e *eventBase[T, V]) EventVersion() uint8 {
  return e.version.EventVersion()
}

func (e *eventBase[T, V]) EventType() EventType {
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

type DomainEvent = eventBase[DomainEventType, NoVersioning]

func NewDomainEvent[V EventVersioning](options ...EventOption[DomainEventType, V]) EventBase[DomainEventType, V] {
  return NewEvent[DomainEventType, V](options...)
}

type DomainEventV1 = eventBase[DomainEventType, V1]

type IntegrationEvent = eventBase[IntegrationEventType, NoVersioning]

func NewIntegrationEvent[V EventVersioning](options ...EventOption[IntegrationEventType, V]) EventBase[IntegrationEventType, V] {
  return NewEvent[IntegrationEventType, V](options...)
}

type IntegrationEventV1 = eventBase[IntegrationEventType, V1]

type IEventBaseConstructable[T EventTyping, V EventVersioning] interface {
  ConstructEventBase(options ...EventOption[T, V])
}

type EventBase[T EventTyping, V EventVersioning] struct {
  eventBase[T, V]
}

func (b *EventBase[T, V]) ConstructEventBase(options ...EventOption[T, V]) {
  b.eventBase = newEvent(options...)
}

func NewEvent[T EventTyping, V EventVersioning](opts ...EventOption[T, V]) EventBase[T, V] {
  return EventBase[T, V]{eventBase: newEvent(opts...)}
}
