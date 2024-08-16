package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/pkg/cqrs"
)

type IVerifyCustomerEmailHandler interface {
  handler.CommandUnit[*VerifyCustomerEmailCommand]
}

func NewVerifyCustomerEmailHandler(parameter cqrs.CommonHandlerParameter) IVerifyCustomerEmailHandler {
  return &verifyCustomerEmail{
    basicHandler: newBasicHandler(&parameter),
  }
}

type verifyCustomerEmail struct {
  basicHandler
}

func (v *verifyCustomerEmail) Handle(ctx context.Context, cmd *VerifyCustomerEmailCommand) status.Object {
  ctx, span := v.tracer.Start(ctx, "verifyCustomerEmail.Handle")
  defer span.End()

  // TODO: Validate token
  userId := types.NullId()

  customers, err := v.repo.FindByIds(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &customers[0]
  ev, err := current.VerifyEmail()
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err = current.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  // Update aggregate
  err = v.repo.Update(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  err = v.publisher.Publish(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Updated()
}
