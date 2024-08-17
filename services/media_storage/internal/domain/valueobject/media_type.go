package valueobject

import (
  sharedErr "github.com/arcorium/rashop/shared/errors"
  "time"
)

type Usage struct {
  Type MediaUsage
  Time time.Time // Used when the usage is UsageTimed
}

func NewMediaUsage(val uint8) (MediaUsage, error) {
  enum := MediaUsage(val)
  if !enum.Valid() {
    return enum, sharedErr.ErrEnumOutOfBounds
  }
  return enum, nil
}

type MediaUsage uint8

const (
  UsageOnce MediaUsage = iota
  //UsageTimed
  UsageFull
)

func (s MediaUsage) Underlying() uint8 {
  return uint8(s)
}

func (s MediaUsage) Valid() bool {
  return s <= UsageFull
}
