package command

import (
  "context"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "rashop/services/customer/pkg/cqrs"
)

type IUpdateCustomerBalanceHandler interface {
  handler.CommandUnit[*UpdateCustomerBalanceCommand]
}

func NewUpdateCustomerBalanceHandler(parameter cqrs.CommonHandlerParameter) IUpdateCustomerBalanceHandler {
  return &updateCustomerBalanceHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type updateCustomerBalanceHandler struct {
  basicHandler
}

func (u *updateCustomerBalanceHandler) Handle(ctx context.Context, cmd *UpdateCustomerBalanceCommand) status.Object {
  ctx, span := u.tracer.Start(ctx, "updateCustomerBalanceHandler.Handle")
  defer span.End()

  customers, err := u.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &customers[0]
  var ev types.Event
  switch cmd.Operator {
  case OperatorMod:
    ev = current.ModifyBalance(cmd.Balance, cmd.Point)
  case OperatorSet:
    ev = current.SetBalance(uint64(cmd.Balance), uint64(cmd.Point))
  }

  err = current.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err = u.repo.Update(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  err = u.publisher.Publish(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Updated()
}
