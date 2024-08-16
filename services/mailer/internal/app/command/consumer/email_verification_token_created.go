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

type IEmailVerificationTokenCreatedHandler interface {
  handler.Consumer[*intev.EmailVerificationTokenCreatedV1]
}

func NewEmailVerificationTokenCreatedHandler(sendHandler command.ISendMailHandler) IEmailVerificationTokenCreatedHandler {
  return &emailVerificationTokenCreatedHandler{
    send:          sendHandler,
    commonHandler: newCommonHandler(),
  }
}

type emailVerificationTokenCreatedHandler struct {
  send command.ISendMailHandler
  commonHandler
}

func (e *emailVerificationTokenCreatedHandler) Handle(ctx context.Context, ev *intev.EmailVerificationTokenCreatedV1) status.Object {
  ctx, span := e.tracer.Start(ctx, "emailVerificationTokenCreated.Handle")
  defer span.End()

  cmd, err := mapper.EmailVerificationToCommand(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  _, stat := e.send.Handle(ctx, &cmd)
  return stat
}
