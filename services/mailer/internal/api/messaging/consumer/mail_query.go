package consumer

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/app/service"
  "github.com/arcorium/rashop/services/mailer/internal/domain/event"
  "github.com/arcorium/rashop/services/mailer/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

func NewMailQueryHandler(svc service.IMailQueryConsumer) IMailQueryHandler {
  return &mailQueryConsumerHandler{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type IMailQueryHandler interface {
  OnMailCreated(ctx context.Context, ev *event.MailCreatedV1) error
  OnMailDeleted(ctx context.Context, ev *event.MailDeletedV1) error
  OnMailStatusUpdated(ctx context.Context, ev *event.MailStatusUpdatedV1) error
}

type mailQueryConsumerHandler struct {
  svc    service.IMailQueryConsumer
  tracer trace.Tracer
}

func (m *mailQueryConsumerHandler) OnMailCreated(ctx context.Context, ev *event.MailCreatedV1) error {
  ctx, span := m.tracer.Start(ctx, "mailQueryConsumerHandler.OnMailCreated")
  defer span.End()

  stat := m.svc.MailCreatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (m *mailQueryConsumerHandler) OnMailDeleted(ctx context.Context, ev *event.MailDeletedV1) error {
  ctx, span := m.tracer.Start(ctx, "mailQueryConsumerHandler.OnMailDeleted")
  defer span.End()

  stat := m.svc.MailDeletedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (m *mailQueryConsumerHandler) OnMailStatusUpdated(ctx context.Context, ev *event.MailStatusUpdatedV1) error {
  ctx, span := m.tracer.Start(ctx, "mailQueryConsumerHandler.OnMailStatusUpdated")
  defer span.End()

  stat := m.svc.MailStatusUpdatedV1(ctx, ev)
  return stat.ErrorWithSpan(span)
}
