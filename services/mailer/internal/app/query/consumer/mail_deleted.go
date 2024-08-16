package consumer

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/domain/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IMailDeletedHandler interface {
  handler.Consumer[*event.MailDeletedV1]
}

func NewMailDeletedHandler(parameter CommonHandlerParameter) IMailDeletedHandler {
  return &mailDeletedHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type mailDeletedHandler struct {
  commonHandler
}

func (m *mailDeletedHandler) Handle(ctx context.Context, ev *event.MailDeletedV1) status.Object {
  ctx, span := m.tracer.Start(ctx, "mailDeletedHandler.Handle")
  defer span.End()

  _, err := m.persistent.Delete(ctx, ev.StartTime, ev.EndTime)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  return status.Deleted()
}
