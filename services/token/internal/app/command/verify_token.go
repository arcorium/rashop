package command

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IVerifyTokenHandler interface {
  handler.CommandUnit[*VerifyTokenCommand]
}

func NewVerifyTokenHandler(parameter CommonHandlerParameter) IVerifyTokenHandler {
  return &verifyTokenHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type verifyTokenHandler struct {
  commonHandler
}

func (v *verifyTokenHandler) Handle(ctx context.Context, cmd *VerifyTokenCommand) status.Object {
  ctx, span := v.tracer.Start(ctx, "verifyTokenHandler.Handle")
  defer span.End()

  // Get aggregate
  token, err := v.persistent.FindByToken(ctx, cmd.Token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  event, err := token.Verify(cmd.Usage)
  err = token.ApplyEvent(event)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Delete in persistent
  err = v.persistent.Delete(ctx, &token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  // Add integration event
  ev := intev.NewTokenVerified(
    token.UserId,
    token.Id,
    token.Token,
    token.Usage.Underlying(),
    token.Type.Underlying(),
    token.ExpiredAt,
    token.GeneratedAt,
  )
  token.AddEvents(ev)

  // Publish message
  return v.Publish(ctx, &token, status.Succeed())
}
