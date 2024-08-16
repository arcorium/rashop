package mapper

import (
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/mailer/constant"
  "github.com/arcorium/rashop/services/mailer/internal/app/command"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
  "github.com/arcorium/rashop/services/mailer/pkg/util"
  "github.com/arcorium/rashop/shared/types"
)

func EmailVerificationToCommand(ev *intev.EmailVerificationTokenCreatedV1) (command.SendMailCommand, error) {
  email, err := types.EmailFromString(ev.Recipient)
  if err != nil {
    return command.SendMailCommand{}, err
  }

  body, bodyType := util.CreateEmailVerificationMailBody(ev.Token, ev.ExpiryTime)

  return command.SendMailCommand{
    Tag:        vob.MailTagEmailVerification,
    Recipients: []types.Email{email},
    Subject:    constant.EmailVerificationSubject,
    BodyType:   bodyType,
    Body:       body,
  }, nil
}

func ResetPasswordToCommand(ev *intev.ResetPasswordTokenCreatedV1) (command.SendMailCommand, error) {
  email, err := types.EmailFromString(ev.Recipient)
  if err != nil {
    return command.SendMailCommand{}, err
  }

  body, bodyType := util.CreateResetPasswordMailBody(ev.Token, ev.ExpiryTime)

  return command.SendMailCommand{
    Tag:        vob.MailTagResetPassword,
    Recipients: []types.Email{email},
    Subject:    constant.ResetPasswordSubject,
    BodyType:   bodyType,
    Body:       body,
  }, nil
}

func LoginTokenToCommand(ev *intev.LoginTokenCreatedV1) (command.SendMailCommand, error) {
  email, err := types.EmailFromString(ev.Recipient)
  if err != nil {
    return command.SendMailCommand{}, err
  }

  body, bodyType := util.CreateLoginTokenMailBody(ev.Token, ev.ExpiryTime)

  return command.SendMailCommand{
    Tag:        vob.MailTagLogin,
    Recipients: []types.Email{email},
    Subject:    constant.LoginSubject,
    BodyType:   bodyType,
    Body:       body,
  }, nil
}
