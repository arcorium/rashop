package mapper

import (
  "github.com/arcorium/rashop/services/mailer/internal/app/dto"
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
)

func ToMailResponse(mail *entity.Mail) dto.MailResponseDTO {
  return dto.MailResponseDTO{
    Id:          mail.Id,
    Status:      mail.Status,
    Tag:         mail.Tag,
    Subject:     mail.Subject,
    Recipients:  mail.Recipients,
    Sender:      mail.Sender,
    SentAt:      mail.SentAt,
    DeliveredAt: mail.DeliveredAt,
  }
}
