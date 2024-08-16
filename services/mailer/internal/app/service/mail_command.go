package service

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/app/command"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
)

type IMailCommand interface {
  Send(ctx context.Context, cmd *command.SendMailCommand) (types.Id, status.Object)
  Delete(ctx context.Context, cmd *command.DeleteMailsCommand) (uint64, status.Object)
}

func NewMailCommand(config MailCommandConfig) IMailCommand {
  return &mailCommandService{
    i: config,
  }
}

func DefaultMailCommandConfig(parameter command.CommonHandlerParameter) MailCommandConfig {
  return MailCommandConfig{
    Delete: command.NewDeleteMailsHandler(parameter),
    Send:   command.NewSendMailHandler(parameter),
  }
}

type MailCommandConfig struct {
  Delete command.IDeleteMailsHandler
  Send   command.ISendMailHandler
}

type mailCommandService struct {
  i MailCommandConfig
}

func (m *mailCommandService) Send(ctx context.Context, cmd *command.SendMailCommand) (types.Id, status.Object) {
  return m.i.Send.Handle(ctx, cmd)
}

func (m *mailCommandService) Delete(ctx context.Context, cmd *command.DeleteMailsCommand) (uint64, status.Object) {
  return m.i.Delete.Handle(ctx, cmd)
}
