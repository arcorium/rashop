package command

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/pkg/cqrs"
)

type IVerificationCustomerEmailRequestHandler interface {
  handler.CommandUnit[*VerificationCustomerEmailRequestCommand]
}

func NewVerificationCustomerEmailRequestHandler(parameter cqrs.CommonHandlerParameter) IVerificationCustomerEmailRequestHandler {
  return &verificationCustomerEmailRequestHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type verificationCustomerEmailRequestHandler struct {
  basicHandler
}

func (f *verificationCustomerEmailRequestHandler) Handle(ctx context.Context, cmd *VerificationCustomerEmailRequestCommand) status.Object {
  ctx, span := f.tracer.Start(ctx, "verificationCustomerEmailRequestHandler.Handle")
  defer span.End()

  // Check user
  customers, err := f.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }
  current := &customers[0]
  current.EmailVerificationRequest()

  //TODO: Forward integration event
  // TODO: Token service should subscribe to those event
  // TODO: Email servic should subscirbe to token_created integrated event
  current.AddEvents(intev.NewCustomerEmailVerificationRequestedV1(current.Id, current.Email, current.Name.User))

  err = f.publisher.Publish(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Succeed()
}
