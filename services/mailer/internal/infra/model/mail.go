package model

import (
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "github.com/arcorium/rashop/shared/util/repo"
  "github.com/uptrace/bun"
  "time"
)

type MailOption = repo.DataAccessModelMapOption[*entity.Mail, *Mail]

func FromMailDomain(ent *entity.Mail, options ...MailOption) Mail {
  mail := Mail{
    Id:          ent.Id.String(),
    Status:      ent.Status.Underlying(),
    Tag:         ent.Tag.Underlying(),
    Subject:     ent.Subject,
    Sender:      ent.Sender.String(),
    DeliveredAt: ent.DeliveredAt,
    CreatedAt:   ent.SentAt,
  }

  for _, recipient := range ent.Recipients {
    mailRecipient := MailRecipient{
      MailId:    mail.Id,
      Recipient: recipient.String(),
    }

    mail.Recipients = append(mail.Recipients, mailRecipient)
  }

  for _, option := range options {
    option(ent, &mail)
  }
  return mail
}

type Mail struct {
  bun.BaseModel `bun:"table:mails,alias:m"`

  Id      string `bun:",nullzero,type:uuid,pk"`
  Status  uint8  `bun:",notnull"`
  Tag     uint8  `bun:",notnull"`
  Subject string `bun:",nullzero,notnull"`
  Sender  string `bun:",nullzero,notnull"`

  DeliveredAt time.Time `bun:",nullzero"`
  CreatedAt   time.Time `bun:",nullzero,notnull"`

  Recipients []MailRecipient `bun:"rel:has-many,join:id=mail_id"`
}

func (m *Mail) ToDomain() (entity.Mail, error) {
  id, err := types.IdFromString(m.Id)
  if err != nil {
    return entity.Mail{}, err
  }

  status, err := vob.NewMailStatus(m.Status)
  if err != nil {
    return entity.Mail{}, err
  }

  tag, err := vob.NewMailTag(m.Tag)
  if err != nil {
    return entity.Mail{}, err
  }

  recipients, ierr := sharedUtil.CastSliceErrsP(m.Recipients, func(from *MailRecipient) (types.Email, error) {
    return types.EmailFromString(from.Recipient)
  })
  if ierr.IsError() {
    return entity.Mail{}, nil
  }

  sender, err := types.EmailFromString(m.Sender)
  if err != nil {
    return entity.Mail{}, err
  }

  return entity.Mail{
    Id:          id,
    Status:      status,
    Tag:         tag,
    Subject:     m.Subject,
    Recipients:  recipients,
    Sender:      sender,
    SentAt:      m.CreatedAt,
    DeliveredAt: m.DeliveredAt,
  }, nil
}
