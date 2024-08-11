package serde

type ISerializer interface {
  Serialize(value any) ([]byte, error)
}
