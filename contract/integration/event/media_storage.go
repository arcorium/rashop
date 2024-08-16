package intev

import (
  "github.com/arcorium/rashop/shared/types"
)

const (
  MediaStorageOneTimeUsedEvent = "media_storage.used"
)

var _ types.Event = (*OneTimeMediaUsedV1)(nil)

func NewOneTimeMediaUsedV1(mediaId types.Id) *OneTimeMediaUsedV1 {
  return &OneTimeMediaUsedV1{
    IntegrationV1: NewV1(),
    MediaId:       mediaId.String(),
  }
}

// OneTimeMediaUsedV1 is message or event when the one time media is already used
type OneTimeMediaUsedV1 struct {
  IntegrationV1
  MediaId string
}

func (c *OneTimeMediaUsedV1) EventName() string {
  return MediaStorageOneTimeUsedEvent
}

func (c *OneTimeMediaUsedV1) Key() (string, bool) {
  return c.MediaId, true
}
