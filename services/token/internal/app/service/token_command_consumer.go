package service

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/token/internal/app/command/consumer"
  "github.com/arcorium/rashop/shared/status"
)

type ITokenCommandConsumer interface {
  CustomerCreated(ctx context.Context, ev *intev.CustomerCreatedV1) status.Object
  CustomerEmailUpdated(ctx context.Context, ev *intev.CustomerEmailUpdatedV1) status.Object
  EmailVerificationRequested(ctx context.Context, ev *intev.CustomerEmailVerificationRequestedV1) status.Object
  LoginByTokenRequested(ctx context.Context, ev *intev.LoginTokenRequestedV1) status.Object
  ResetPasswordRequested(ctx context.Context, ev *intev.CustomerResetPasswordRequestedV1) status.Object
}

func NewTokenCommandConsumer(config TokenCommandConsumerFactory) ITokenCommandConsumer {
  return &tokenCommandConsumerService{
    i: config,
  }
}

func DefaultTokenCommandConsumerFactory(parameter consumer.CommonHandlerParameter) TokenCommandConsumerFactory {
  return TokenCommandConsumerFactory{
    CustomerCreated:            consumer.NewCustomerCreatedHandler(parameter),
    CustomerEmailUpdated:       consumer.NewCustomerEmailChangedHandler(parameter),
    EmailVerificationRequested: consumer.NewEmailVerificationRequestedHandler(parameter),
    LoginByTokenRequested:      consumer.NewLoginRequestedHandler(parameter),
    ResetPasswordRequested:     consumer.NewResetPasswordRequestedHandler(parameter),
  }
}

type TokenCommandConsumerFactory struct {
  CustomerCreated            consumer.ICustomerCreatedHandler
  CustomerEmailUpdated       consumer.ICustomerEmailChangedHandler
  EmailVerificationRequested consumer.IEmailVerificationRequestedHandler
  LoginByTokenRequested      consumer.ILoginRequestedHandler
  ResetPasswordRequested     consumer.IResetPasswordRequestedHandler
}

type tokenCommandConsumerService struct {
  i TokenCommandConsumerFactory
}

func (t *tokenCommandConsumerService) CustomerCreated(ctx context.Context, ev *intev.CustomerCreatedV1) status.Object {
  return t.i.CustomerCreated.Handle(ctx, ev)
}

func (t *tokenCommandConsumerService) CustomerEmailUpdated(ctx context.Context, ev *intev.CustomerEmailUpdatedV1) status.Object {
  return t.i.CustomerEmailUpdated.Handle(ctx, ev)
}

func (t *tokenCommandConsumerService) EmailVerificationRequested(ctx context.Context, ev *intev.CustomerEmailVerificationRequestedV1) status.Object {
  return t.i.EmailVerificationRequested.Handle(ctx, ev)
}

func (t *tokenCommandConsumerService) LoginByTokenRequested(ctx context.Context, ev *intev.LoginTokenRequestedV1) status.Object {
  return t.i.LoginByTokenRequested.Handle(ctx, ev)
}

func (t *tokenCommandConsumerService) ResetPasswordRequested(ctx context.Context, ev *intev.CustomerResetPasswordRequestedV1) status.Object {
  return t.i.ResetPasswordRequested.Handle(ctx, ev)
}
