package command

import (
  "github.com/arcorium/rashop/services/mailer/constant"
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type SendMailCommand struct {
  Tag                vob.MailTag
  Recipients         []types.Email
  Sender             types.NullableEmail
  Subject            string
  BodyType           vob.BodyType
  Body               string
  EmbeddedMediaIds   []types.Id
  AttachmentMediaIds []types.Id
}

func (m *SendMailCommand) ToDomain() (entity.Mail, error) {
  id, err := types.NewId()
  if err != nil {
    return entity.Mail{}, err
  }

  return entity.Mail{
    Id:         id,
    Status:     vob.MailStatusPending,
    Tag:        m.Tag,
    Subject:    m.Subject,
    Body:       m.Body,
    BodyType:   m.BodyType,
    Recipients: m.Recipients,
    Sender:     m.Sender.ValueOr(constant.DEFAULT_EMAIL_SENDER),
    SentAt:     time.Now(),
  }, nil
}

type DeleteMailsCommand struct {
  StartTime time.Time
  EndTime   time.Time
}
