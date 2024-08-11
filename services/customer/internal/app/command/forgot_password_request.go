package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
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

  err = f.publisher.PublishAggregate(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Success()
}
