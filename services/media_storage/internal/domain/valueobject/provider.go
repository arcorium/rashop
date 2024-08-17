package valueobject

import sharedErr "github.com/arcorium/rashop/shared/errors"

func NewStorageProvider(val uint8) (StorageProvider, error) {
  enum := StorageProvider(val)
  if !enum.Valid() {
    return enum, sharedErr.ErrEnumOutOfBounds
  }
  return enum, nil
}

type StorageProvider uint8

const (
  ProviderMinIO StorageProvider = iota
  ProviderUnknown
)

func (s StorageProvider) Underlying() uint8 {
  return uint8(s)
}

func (s StorageProvider) Valid() bool {
  return s < ProviderUnknown
}
