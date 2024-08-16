package consumer

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
  "github.com/arcorium/rashop/services/mailer/internal/domain/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IMailCreatedHandler interface {
  handler.Consumer[*event.MailCreatedV1]
}

func NewMailCreatedHandler(parameter CommonHandlerParameter) IMailCreatedHandler {
  return &mailCreatedHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type mailCreatedHandler struct {
  commonHandler
}

func (m *mailCreatedHandler) Handle(ctx context.Context, ev *event.MailCreatedV1) status.Object {
  ctx, span := m.tracer.Start(ctx, "mailCreatedHandler.Handle")
  defer span.End()

  mail := entity.Mail{}
  err := mail.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  err = m.persistent.Create(ctx, &mail)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  return status.Created()
}
