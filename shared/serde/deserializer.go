package serde

type IDeserializer[T any] interface {
  Deserialize([]byte) (*T, error)
}

type IDeserializerAny interface {
  Deserialize([]byte, any) error
}
