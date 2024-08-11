package types

import (
  "context"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/serde"
  "strconv"
)

const (
  METADATA_IDENTITY_KEY   = "id"
  METADATA_VERSION_KEY    = "version"
  METADATA_EVENT_TYPE_KEY = "type"
)

type EventMetadataCtxKey struct{}

type Record = KeyVal[string, string]

type Metadata []Record

func (m Metadata) ToKafkaRecordHeader() []sarama.RecordHeader {
  result := make([]sarama.RecordHeader, len(m))
  for _, record := range m {
    result = append(result, sarama.RecordHeader{
      Key:   []byte(record.Key),
      Value: []byte(record.Val),
    })
  }
  return result
}

func (m Metadata) Serialize(serializer serde.ISerializer) ([]byte, error) {
  return serializer.Serialize(m)
}

type eventMetadataBuilder[T any] interface {
  Handle(T) *EventMetadata
}

type metadataEventBuilder struct{}

func (m metadataEventBuilder) Handle(metadata Metadata) *EventMetadata {
  result := &EventMetadata{}
  for _, record := range metadata {
    switch record.Key {
    case METADATA_IDENTITY_KEY:
      result.Id = record.Val
    case METADATA_EVENT_TYPE_KEY:
      result.Type = NewEventType(record.Val)
    case METADATA_VERSION_KEY:
      result.Version = uint8(Must(strconv.ParseUint(record.Val, 10, 8)))
    }
  }
  return result
}

type kafkaRecordEventBuilder struct{}

func (m kafkaRecordEventBuilder) Handle(headers []*sarama.RecordHeader) *EventMetadata {
  result := &EventMetadata{}
  for _, record := range headers {
    switch string(record.Key) {
    case METADATA_IDENTITY_KEY:
      result.Id = string(record.Value)
    case METADATA_EVENT_TYPE_KEY:
      result.Type = NewEventType(string(record.Value))
    case METADATA_VERSION_KEY:
      result.Version = uint8(Must(strconv.ParseUint(string(record.Value), 10, 8)))
    }
  }
  return result
}

type contextEventBuilder struct{}

func (m contextEventBuilder) Handle(ctx context.Context) *EventMetadata {
  value := ctx.Value(EventMetadataCtxKey{})
  if value == nil {
    return nil
  }
  metadata, ok := value.(*EventMetadata)
  if !ok {
    return nil
  }
  return metadata
}

func FromMetadata() eventMetadataBuilder[Metadata] {
  return metadataEventBuilder{}
}

func FromKafkaRecords() eventMetadataBuilder[[]*sarama.RecordHeader] {
  return kafkaRecordEventBuilder{}
}

func FromContext() eventMetadataBuilder[context.Context] {
  return contextEventBuilder{}
}

func NewEventMetadata[T any](val T, builder eventMetadataBuilder[T]) *EventMetadata {
  return builder.Handle(val)
}

type EventMetadata struct {
  Id      string
  Version uint8
  Type    EventType
}

func (m *EventMetadata) AsContext(parent context.Context) context.Context {
  return context.WithValue(parent, EventMetadataCtxKey{}, m)
}

type EventVersioning interface {
  EventVersion() uint8
}

func NewEventType(val string) EventType {
  switch val {
  case string(EventTypeUnknown):
    return EventTypeUnknown
  case string(EventTypeDomain):
    return EventTypeDomain
  case string(EventTypeIntegration):
    return EventTypeIntegration
  }
  return EventTypeUnknown
}

type EventType string

const (
  EventTypeUnknown     EventType = "unknown"
  EventTypeDomain      EventType = "domain"
  EventTypeIntegration EventType = "integration"
)

type EventTyping interface {
  EventType() EventType
}

type NoEventTyping struct{}

func (NoEventTyping) EvenType() EventType {
  return EventTypeUnknown
}

type NoVersioning struct{}

func (e NoVersioning) EventVersion() uint8 {
  return 0
}

type V1 struct{}

func (e V1) EventVersion() uint8 {
  return 1
}

type V2 struct{}

func (e V2) EventVersion() uint8 {
  return 2
}
