package command

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/pkg/cqrs"
)

type IForgotCustomerPasswordRequestHandler interface {
  handler.CommandUnit[*ForgotCustomerPasswordRequestCommand]
}

func NewForgotCustomerPasswordRequestHandler(parameter cqrs.CommonHandlerParameter) IForgotCustomerPasswordRequestHandler {
  return &forgotCustomerPasswordRequestHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type forgotCustomerPasswordRequestHandler struct {
  basicHandler
}

func (f *forgotCustomerPasswordRequestHandler) Handle(ctx context.Context, cmd *ForgotCustomerPasswordRequestCommand) status.Object {
  ctx, span := f.tracer.Start(ctx, "forgotCustomerPasswordRequestHandler.Handle")
  defer span.End()

  // Check user
  customers, err := f.repo.FindByEmails(ctx, cmd.Email)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &customers[0]
  current.ForgotPasswordRequest()

  // add integration event
  current.AddEvents(intev.NewCustomerResetPasswordRequestedV1(current.Id, current.Email, current.Name.User))

  err = f.publisher.Publish(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Succeed()
}
