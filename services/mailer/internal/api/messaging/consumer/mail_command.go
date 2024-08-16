package consumer

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/mailer/internal/app/service"
  "github.com/arcorium/rashop/services/mailer/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

func NewMailCommandHandler(svc service.IMailCommandConsumer) IMailCommandHandler {
  return &mailCommandConsumerHandler{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type IMailCommandHandler interface {
  OnEmailVerificationTokenCreated(ctx context.Context, ev *intev.EmailVerificationTokenCreatedV1) error
  OnResetPasswordTokenCreated(ctx context.Context, ev *intev.ResetPasswordTokenCreatedV1) error
  OnLoginTokenCreated(ctx context.Context, ev *intev.LoginTokenCreatedV1) error
}

type mailCommandConsumerHandler struct {
  svc    service.IMailCommandConsumer
  tracer trace.Tracer
}

func (m *mailCommandConsumerHandler) OnEmailVerificationTokenCreated(ctx context.Context, ev *intev.EmailVerificationTokenCreatedV1) error {
  ctx, span := m.tracer.Start(ctx, "mailCommandConsumerHandler.OnEmailVerificationTokenCreated")
  defer span.End()

  stat := m.svc.EmailVerificationTokenCreated(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (m *mailCommandConsumerHandler) OnResetPasswordTokenCreated(ctx context.Context, ev *intev.ResetPasswordTokenCreatedV1) error {
  ctx, span := m.tracer.Start(ctx, "mailCommandConsumerHandler.OnResetPasswordTokenCreated")
  defer span.End()

  stat := m.svc.ResetPasswordTokenCreated(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (m *mailCommandConsumerHandler) OnLoginTokenCreated(ctx context.Context, ev *intev.LoginTokenCreatedV1) error {
  ctx, span := m.tracer.Start(ctx, "mailCommandConsumerHandler.OnLoginTokenCreated")
  defer span.End()

  stat := m.svc.LoginTokenCreated(ctx, ev)
  return stat.ErrorWithSpan(span)
}
