package valueobject

import sharedErr "github.com/arcorium/rashop/shared/errors"

func NewUsage(val uint8) (TokenUsage, error) {
  enum := TokenUsage(val)
  if !enum.Valid() {
    return 0, sharedErr.ErrEnumOutOfBounds
  }
  return enum, nil
}

type TokenUsage uint8

const (
  UsageEmailVerification TokenUsage = iota
  UsageResetPassword
  UsageLogin
  UsageGeneral
)

func (t TokenUsage) Underlying() uint8 {
  return uint8(t)
}

func (t TokenUsage) Valid() bool {
  return t <= UsageGeneral
}
