package consumer

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/mailer/internal/app/command"
  "github.com/arcorium/rashop/services/mailer/internal/app/mapper"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IResetPasswordTokenCreatedHandler interface {
  handler.Consumer[*intev.ResetPasswordTokenCreatedV1]
}

func NewResetPasswordTokenCreatedHandler(sendHandler command.ISendMailHandler) IResetPasswordTokenCreatedHandler {
  return &resetPasswordTokenCreatedHandler{
    send:          sendHandler,
    commonHandler: newCommonHandler(),
  }
}

type resetPasswordTokenCreatedHandler struct {
  send command.ISendMailHandler
  commonHandler
}

func (e *resetPasswordTokenCreatedHandler) Handle(ctx context.Context, ev *intev.ResetPasswordTokenCreatedV1) status.Object {
  ctx, span := e.tracer.Start(ctx, "resetPasswordTokenCreated.Handle")
  defer span.End()

  cmd, err := mapper.ResetPasswordToCommand(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  _, stat := e.send.Handle(ctx, &cmd)
  return stat
}
