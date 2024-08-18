package consumer

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/token/internal/app/service"
  "github.com/arcorium/rashop/services/token/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

func NewTokenCommand(svc service.ITokenCommandConsumer) ITokenCommandHandler {
  return &tokenCommandConsumerHandler{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type ITokenCommandHandler interface {
  OnCustomerCreated(ctx context.Context, ev *intev.CustomerCreatedV1) error
  OnCustomerEmailUpdated(ctx context.Context, ev *intev.CustomerEmailUpdatedV1) error
  OnEmailVerificationRequested(ctx context.Context, ev *intev.CustomerEmailVerificationRequestedV1) error
  OnLoginRequested(ctx context.Context, ev *intev.LoginTokenRequestedV1) error
  OnResetPasswordRequested(ctx context.Context, ev *intev.CustomerResetPasswordRequestedV1) error
}

type tokenCommandConsumerHandler struct {
  svc    service.ITokenCommandConsumer
  tracer trace.Tracer
}

func (t *tokenCommandConsumerHandler) OnCustomerCreated(ctx context.Context, ev *intev.CustomerCreatedV1) error {
  ctx, span := t.tracer.Start(ctx, "tokenCommandConsumerHandler.OnCustomerCreated")
  defer span.End()

  stat := t.svc.CustomerCreated(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (t *tokenCommandConsumerHandler) OnCustomerEmailUpdated(ctx context.Context, ev *intev.CustomerEmailUpdatedV1) error {
  ctx, span := t.tracer.Start(ctx, "tokenCommandConsumerHandler.OnCustomerEmailUpdated")
  defer span.End()

  stat := t.svc.CustomerEmailUpdated(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (t *tokenCommandConsumerHandler) OnEmailVerificationRequested(ctx context.Context, ev *intev.CustomerEmailVerificationRequestedV1) error {
  ctx, span := t.tracer.Start(ctx, "tokenCommandConsumerHandler.OnEmailVerificationRequested")
  defer span.End()

  stat := t.svc.EmailVerificationRequested(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (t *tokenCommandConsumerHandler) OnLoginRequested(ctx context.Context, ev *intev.LoginTokenRequestedV1) error {
  ctx, span := t.tracer.Start(ctx, "tokenCommandConsumerHandler.OnLoginRequested")
  defer span.End()

  stat := t.svc.LoginByTokenRequested(ctx, ev)
  return stat.ErrorWithSpan(span)
}

func (t *tokenCommandConsumerHandler) OnResetPasswordRequested(ctx context.Context, ev *intev.CustomerResetPasswordRequestedV1) error {
  ctx, span := t.tracer.Start(ctx, "tokenCommandConsumerHandler.OnResetPasswordRequested")
  defer span.End()

  stat := t.svc.ResetPasswordRequested(ctx, ev)
  return stat.ErrorWithSpan(span)
}
