package entity

import (
  "errors"
  "github.com/arcorium/rashop/services/mailer/internal/domain/event"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "time"
)

var (
  ErrUnknownEvent         = errors.New("unknown event")
  ErrUpdateNonPendingMail = errors.New("couldn't update non-pending mail")
)

func CreateMail(mail *Mail) (Mail, types.Event) {
  result := Mail{}
  ev := &event.MailCreatedV1{
    DomainV1:   event.NewV1(),
    MailId:     mail.Id.String(),
    Tag:        mail.Tag.Underlying(),
    Subject:    mail.Subject,
    BodyType:   mail.BodyType.Underlying(),
    Recipients: sharedUtil.CastSlice(mail.Recipients, sharedUtil.ToString[types.Email]),
    Sender:     mail.Sender.String(),
    SentAt:     mail.SentAt,
  }

  mail.AddEvents(ev)
  return result, ev
}

var _ types.Aggregate = (*Mail)(nil)

type Mail struct {
  types.AggregateBase
  types.AggregateHelper
  Id          types.Id
  Status      vob.MailStatus
  Tag         vob.MailTag
  Subject     string
  Body        string
  BodyType    vob.BodyType
  Recipients  []types.Email
  Sender      types.Email
  SentAt      time.Time
  DeliveredAt time.Time
}

func (m *Mail) Failed() (types.Event, error) {
  if m.Status != vob.MailStatusPending {
    return nil, ErrUpdateNonPendingMail
  }

  ev := &event.MailStatusUpdatedV1{
    DomainV1: event.NewV1(),
    MailId:   m.Id.String(),
    Status:   vob.MailStatusFailed.Underlying(),
  }
  m.AddEvents(ev)
  return ev, nil
}

func (m *Mail) Delivered() (types.Event, error) {
  if m.Status != vob.MailStatusPending {
    return nil, ErrUpdateNonPendingMail
  }

  ev := &event.MailStatusUpdatedV1{
    DomainV1: event.NewV1(),
    MailId:   m.Id.String(),
    Status:   vob.MailStatusDelivered.Underlying(),
  }
  m.AddEvents(ev)
  return ev, nil
}

func (m *Mail) Identity() string {
  return m.Id.String()
}

func (m *Mail) ApplyEvent(ev types.Event) error {
  switch cur := ev.(type) {
  case *event.MailCreatedV1:
    recipients, _ := sharedUtil.CastSliceErrs(cur.Recipients, types.EmailFromString)
    m.Id = types.Must(types.IdFromString(cur.MailId))
    m.Tag = types.Must(vob.NewMailTag(cur.Tag))
    m.Subject = cur.Subject
    m.BodyType = types.Must(vob.NewBodyType(cur.BodyType))
    m.Recipients = recipients
    m.Sender = types.Must(types.EmailFromString(cur.Sender))
    m.SentAt = cur.SentAt
    m.MarkCreated()
  case *event.MailStatusUpdatedV1:
    m.Status = types.Must(vob.NewMailStatus(cur.Status))
    m.MarkUpdated()
  default:
    return ErrUnknownEvent
  }
  return nil
}
