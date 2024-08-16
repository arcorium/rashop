package service

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/mailer/internal/app/command"
  commandCon "github.com/arcorium/rashop/services/mailer/internal/app/command/consumer"
  "github.com/arcorium/rashop/shared/status"
)

type IMailCommandConsumer interface {
  EmailVerificationTokenCreated(ctx context.Context, ev *intev.EmailVerificationTokenCreatedV1) status.Object
  ResetPasswordTokenCreated(ctx context.Context, ev *intev.ResetPasswordTokenCreatedV1) status.Object
  LoginTokenCreated(ctx context.Context, ev *intev.LoginTokenCreatedV1) status.Object
}

func NewMailCommandConsumer(config MailCommandConsumerConfig) IMailCommandConsumer {
  return &mailConsumerService{
    i: config,
  }
}

func DefaultMailCommandConsumerConfig(sendHandler command.ISendMailHandler) MailCommandConsumerConfig {
  return MailCommandConsumerConfig{
    EmailVerificationTokenCreated: commandCon.NewEmailVerificationTokenCreatedHandler(sendHandler),
    ResetPasswordTokenCreated:     commandCon.NewResetPasswordTokenCreatedHandler(sendHandler),
    LoginTokenCreated:             commandCon.NewLoginTokenCreatedHandler(sendHandler),
  }
}

type MailCommandConsumerConfig struct {
  EmailVerificationTokenCreated commandCon.IEmailVerificationTokenCreatedHandler
  ResetPasswordTokenCreated     commandCon.IResetPasswordTokenCreatedHandler
  LoginTokenCreated             commandCon.ILoginTokenCreatedHandler
}

type mailConsumerService struct {
  i MailCommandConsumerConfig
}

func (m *mailConsumerService) EmailVerificationTokenCreated(ctx context.Context, ev *intev.EmailVerificationTokenCreatedV1) status.Object {
  return m.i.EmailVerificationTokenCreated.Handle(ctx, ev)
}

func (m *mailConsumerService) ResetPasswordTokenCreated(ctx context.Context, ev *intev.ResetPasswordTokenCreatedV1) status.Object {
  return m.i.ResetPasswordTokenCreated.Handle(ctx, ev)
}

func (m *mailConsumerService) LoginTokenCreated(ctx context.Context, ev *intev.LoginTokenCreatedV1) status.Object {
  return m.i.LoginTokenCreated.Handle(ctx, ev)
}
