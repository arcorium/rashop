package intev

import (
  "github.com/arcorium/rashop/shared/types"
  "time"
)

const (
  MailerDeliveredEvent = "mailer.delivered"
)

var _ types.Event = (*MailDeliveredV1)(nil)

func NewMailDeliveredV1(mailId types.Id, succeed bool, deliverTime time.Time) *MailDeliveredV1 {
  return &MailDeliveredV1{
    IntegrationV1: NewV1(),
    MailId:        mailId.String(),
    Succeed:       succeed,
    DeliveredAt:   deliverTime,
  }
}

type MailDeliveredV1 struct {
  IntegrationV1
  MailId      string
  Succeed     bool
  DeliveredAt time.Time
}

func (c *MailDeliveredV1) EventName() string {
  return MailerDeliveredEvent
}

func (c *MailDeliveredV1) Key() (string, bool) {
  return c.MailId, true
}
