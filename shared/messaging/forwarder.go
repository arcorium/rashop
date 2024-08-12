package messaging

import (
  "context"
)

const HEADER_ERROR_KEY = "error"

type IForwarder[T any] interface {
  Forward(ctx context.Context, message T, err error) error
}
