package service

import (
  "context"
  queryCon "github.com/arcorium/rashop/services/mailer/internal/app/query/consumer"
  "github.com/arcorium/rashop/services/mailer/internal/domain/event"
  "github.com/arcorium/rashop/shared/status"
)

type IMailQueryConsumer interface {
  MailCreatedV1(ctx context.Context, ev *event.MailCreatedV1) status.Object
  MailDeletedV1(ctx context.Context, ev *event.MailDeletedV1) status.Object
  MailStatusUpdatedV1(ctx context.Context, ev *event.MailStatusUpdatedV1) status.Object
}

func NewMailQueryConsumer(config MailQueryConsumerConfig) IMailQueryConsumer {
  return &mailQueryConsumer{
    i: config,
  }
}

func DefaultMailQueryConsumerConfig(parameter queryCon.CommonHandlerParameter) MailQueryConsumerConfig {
  return MailQueryConsumerConfig{
    MailCreated:       queryCon.NewMailCreatedHandler(parameter),
    MailDeleted:       queryCon.NewMailDeletedHandler(parameter),
    MailStatusUpdated: queryCon.NewMailStatusUpdatedHandler(parameter),
  }
}

type MailQueryConsumerConfig struct {
  MailCreated       queryCon.IMailCreatedHandler
  MailDeleted       queryCon.IMailDeletedHandler
  MailStatusUpdated queryCon.IMailStatusUpdatedHandler
}

type mailQueryConsumer struct {
  i MailQueryConsumerConfig
}

func (m *mailQueryConsumer) MailCreatedV1(ctx context.Context, ev *event.MailCreatedV1) status.Object {
  return m.i.MailCreated.Handle(ctx, ev)
}

func (m *mailQueryConsumer) MailDeletedV1(ctx context.Context, ev *event.MailDeletedV1) status.Object {
  return m.i.MailDeleted.Handle(ctx, ev)
}

func (m *mailQueryConsumer) MailStatusUpdatedV1(ctx context.Context, ev *event.MailStatusUpdatedV1) status.Object {
  return m.i.MailStatusUpdated.Handle(ctx, ev)
}
