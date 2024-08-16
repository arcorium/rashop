package valueobject

import (
  sharedErr "github.com/arcorium/rashop/shared/errors"
  "math"
)

func NewBodyType(val uint8) (BodyType, error) {
  enum := BodyType(val)
  if !enum.Valid() {
    return math.MaxUint8, sharedErr.ErrEnumOutOfBounds
  }
  return enum, nil
}

type BodyType uint8

const (
  BodyTypePlain BodyType = iota
  BodyTypeHtml
)

func (b BodyType) Underlying() uint8 {
  return uint8(b)
}

func (b BodyType) Valid() bool {
  return b <= BodyTypeHtml
}

func (b BodyType) String() string {
  switch b {
  case BodyTypePlain:
    return "text/plain"
  case BodyTypeHtml:
    return "text/html"
  }
  return "text/plain" // WARN: Panic?
}

func NewMailStatus(val uint8) (MailStatus, error) {
  enum := MailStatus(val)
  if !enum.Valid() {
    return math.MaxUint8, sharedErr.ErrEnumOutOfBounds
  }
  return enum, nil
}

type MailStatus uint8

const (
  MailStatusPending MailStatus = iota
  MailStatusFailed
  MailStatusDelivered
)

func (m MailStatus) Underlying() uint8 {
  return uint8(m)
}

func (m MailStatus) Valid() bool {
  return m <= MailStatusDelivered
}

func (m MailStatus) String() string {
  switch m {
  case MailStatusPending:
    return "pending"
  case MailStatusFailed:
    return "failed"
  case MailStatusDelivered:
    return "delivered"
  }
  return "unknown"
}

func NewMailTag(val uint8) (MailTag, error) {
  enum := MailTag(val)
  if !enum.Valid() {
    return math.MaxUint8, sharedErr.ErrEnumOutOfBounds
  }
  return enum, nil
}

type MailTag uint8

const (
  MailTagEmailVerification MailTag = iota
  MailTagResetPassword
  MailTagLogin
  MailTagOther
)

func (m MailTag) Underlying() uint8 {
  return uint8(m)
}

func (m MailTag) Valid() bool {
  return m <= MailTagOther
}

func (m MailTag) String() string {
  switch m {
  case MailTagEmailVerification:
    return "verification"
  case MailTagResetPassword:
    return "reset-password"
  case MailTagLogin:
    return "login"
  case MailTagOther:
    return "other"
  }
  return "unknown"
}
