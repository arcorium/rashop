package mapper

import (
  mailerv1 "github.com/arcorium/rashop/proto/gen/go/mailer/v1"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
)

func toProtoTag(tag vob.MailTag) mailerv1.Tag {
  switch tag {
  case vob.MailTagEmailVerification:
    return mailerv1.Tag_EmailVerification
  case vob.MailTagResetPassword:
    return mailerv1.Tag_ResetPassword
  case vob.MailTagLogin:
    return mailerv1.Tag_Login
  case vob.MailTagOther:
    return mailerv1.Tag_Other
  }
  panic("unknown tag")
}

func toProtoStatus(status vob.MailStatus) mailerv1.Status {
  switch status {
  case vob.MailStatusPending:
    return mailerv1.Status_Pending
  case vob.MailStatusFailed:
    return mailerv1.Status_Failed
  case vob.MailStatusDelivered:
    return mailerv1.Status_Delivered
  }
  panic("unknown status")
}
