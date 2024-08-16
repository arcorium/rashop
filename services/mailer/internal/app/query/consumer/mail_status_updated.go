package consumer

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/domain/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IMailStatusUpdatedHandler interface {
  handler.Consumer[*event.MailStatusUpdatedV1]
}

func NewMailStatusUpdatedHandler(parameter CommonHandlerParameter) IMailStatusUpdatedHandler {
  return &mailStatusUpdatedHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type mailStatusUpdatedHandler struct {
  commonHandler
}

func (m *mailStatusUpdatedHandler) Handle(ctx context.Context, ev *event.MailStatusUpdatedV1) status.Object {
  ctx, span := m.tracer.Start(ctx, "mailStatusUpdatedHandler.Handle")
  defer span.End()

  mailId, err := types.IdFromString(ev.MailId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  mails, err := m.persistent.FindByIds(ctx, mailId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &mails[0]

  err = current.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  err = m.persistent.Update(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  return status.Updated()
}
