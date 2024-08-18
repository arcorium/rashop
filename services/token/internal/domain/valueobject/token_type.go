package valueobject

import sharedErr "github.com/arcorium/rashop/shared/errors"

func NewType(val uint8) (TokenType, error) {
  enum := TokenType(val)
  if !enum.Valid() {
    return 0, sharedErr.ErrEnumOutOfBounds
  }
  return enum, nil
}

type TokenType uint8

const (
  TypeString TokenType = iota
  TypePIN
  TypeAlphanumericPIN
)

func (t TokenType) Underlying() uint8 {
  return uint8(t)
}

func (t TokenType) Valid() bool {
  return t <= TypeAlphanumericPIN
}
