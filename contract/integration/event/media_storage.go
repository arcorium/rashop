package intev

import (
  "github.com/arcorium/rashop/contract/integration/enum"
  "github.com/arcorium/rashop/shared/types"
)

const (
  MediaStorageOneTimeUsedEvent = "media.used"
  MediaStoredEvent             = "media.stored"
  MediaDeletedEvent            = "media.deleted"
)

// TODO: Move this event into email
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

var _ types.Event = (*MediaStoredV1)(nil)

func NewMediaStoredV1(userId, mediaId types.Id, usage uint8) *MediaStoredV1 {
  return &MediaStoredV1{
    IntegrationV1: NewV1(),
    UserId:        userId.String(),
    MediaId:       mediaId.String(),
    Usage:         enum.MediaUsage(usage),
  }
}

// MediaStoredV1 is message or event when the one time media is already used
type MediaStoredV1 struct {
  IntegrationV1
  UserId  string
  MediaId string
  Usage   enum.MediaUsage
}

func (c *MediaStoredV1) EventName() string {
  return MediaStoredEvent
}

func (c *MediaStoredV1) Key() (string, bool) {
  return c.MediaId, true
}

var _ types.Event = (*MediaDeletedV1)(nil)

func NewMediaDeletedV1(mediaId types.Id, usage uint8) *MediaDeletedV1 {
  return &MediaDeletedV1{
    IntegrationV1: NewV1(),
    MediaId:       mediaId.String(),
    Usage:         enum.MediaUsage(usage),
  }
}

// MediaDeletedV1 is message or event when the one time media is already used
type MediaDeletedV1 struct {
  IntegrationV1
  MediaId string
  Usage   enum.MediaUsage
}

func (c *MediaDeletedV1) EventName() string {
  return MediaDeletedEvent
}

func (c *MediaDeletedV1) Key() (string, bool) {
  return c.MediaId, true
}
