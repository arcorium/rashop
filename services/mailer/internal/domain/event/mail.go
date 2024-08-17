package event

import (
  "github.com/arcorium/rashop/shared/types"
  "time"
)

const (
  MailerCreatedEvent       = "mail.created"
  MailerStatusUpdatedEvent = "mail.updated"
  MailerDeletedEvent       = "mail.deleted"
)

var _ types.Event = (*MailCreatedV1)(nil)

type MailCreatedV1 struct {
  DomainV1
  MailId     string
  Tag        uint8
  Subject    string
  BodyType   uint8
  Recipients []string
  Sender     string
  SentAt     time.Time
}

func (c *MailCreatedV1) EventName() string {
  return MailerCreatedEvent
}

func (c *MailCreatedV1) Key() (string, bool) {
  return c.MailId, true
}

var _ types.Event = (*MailStatusUpdatedV1)(nil)

type MailStatusUpdatedV1 struct {
  DomainV1
  MailId string
  Status uint8
}

func (c *MailStatusUpdatedV1) EventName() string {
  return MailerStatusUpdatedEvent
}

func (c *MailStatusUpdatedV1) Key() (string, bool) {
  return c.MailId, true
}

var _ types.Event = (*MailDeletedV1)(nil)

type MailDeletedV1 struct {
  DomainV1
  StartTime time.Time
  EndTime   time.Time
}

func (c *MailDeletedV1) EventName() string {
  return MailerDeletedEvent
}

//func (c *MailDeletedV1) Key() (string, bool) {
//  return c.MailId, true
//}
