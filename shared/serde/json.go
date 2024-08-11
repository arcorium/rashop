package serde

import "encoding/json"

type JsonSerializer struct{}

func (j JsonSerializer) Serialize(value any) ([]byte, error) {
  return json.Marshal(value)
}

type JsonDeserializer[T any] struct{}

func (j JsonDeserializer[T]) Deserialize(bytes []byte) (*T, error) {
  var t *T
  err := json.Unmarshal(bytes, &t)
  return t, err
}

type JsonAnyDeserializer struct{}

func (j JsonAnyDeserializer) Deserialize(bytes []byte, result any) error {
  return json.Unmarshal(bytes, result)
}

func ToJSON(value any) ([]byte, error) {
  return JsonSerializer{}.Serialize(value)
}

func FromJSON[T any](value []byte) (*T, error) {
  return JsonDeserializer[T]{}.Deserialize(value)
}

func BinJSON[T any](value []byte, val T) error {
  return JsonAnyDeserializer{}.Deserialize(value, val)
}
