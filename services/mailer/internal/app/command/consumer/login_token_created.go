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

type ILoginTokenCreatedHandler interface {
  handler.Consumer[*intev.LoginTokenCreatedV1]
}

func NewLoginTokenCreatedHandler(sendHandler command.ISendMailHandler) ILoginTokenCreatedHandler {
  return &loginTokenCreatedHandler{
    send:          sendHandler,
    commonHandler: newCommonHandler(),
  }
}

type loginTokenCreatedHandler struct {
  send command.ISendMailHandler
  commonHandler
}

func (e *loginTokenCreatedHandler) Handle(ctx context.Context, ev *intev.LoginTokenCreatedV1) status.Object {
  ctx, span := e.tracer.Start(ctx, "loginTokenCreated.Handle")
  defer span.End()

  cmd, err := mapper.LoginTokenToCommand(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  _, stat := e.send.Handle(ctx, &cmd)
  return stat
}
