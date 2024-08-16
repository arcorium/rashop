package repository

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
)

type IMailClient interface {
  Send(ctx context.Context, mail *entity.Mail, attachments []vob.Attachment) error
  Close(ctx context.Context) error
}
