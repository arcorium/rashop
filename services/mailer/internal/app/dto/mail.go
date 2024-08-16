package dto

import (
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type MailResponseDTO struct {
  Id          types.Id
  Status      vob.MailStatus
  Tag         vob.MailTag
  Subject     string
  Recipients  []types.Email
  Sender      types.Email
  SentAt      time.Time
  DeliveredAt time.Time
}
