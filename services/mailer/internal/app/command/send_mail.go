package command

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/mailer/internal/domain/entity"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "time"
)

type ISendMailHandler interface {
  handler.Command[*SendMailCommand, types.Id]
}

func NewSendMailHandler(parameter CommonHandlerParameter) ISendMailHandler {
  return &sendMailHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type sendMailHandler struct {
  commonHandler
}

func (s *sendMailHandler) getAttachments(ctx context.Context, cmd *SendMailCommand) ([]vob.Attachment, []types.Event) {
  var attachments []vob.Attachment
  var events []types.Event

  for _, mediaId := range cmd.EmbeddedMediaIds {
    media, err := s.mediaClient.FindById(ctx, mediaId)
    if err != nil {
      continue
    }
    attachments = append(attachments, vob.Attachment{
      IsEmbedded: true,
      Filename:   media.Name,
      Data:       media.Data,
    })

    if media.Type == vob.MediaTypeOneTime {
      events = append(events, intev.NewOneTimeMediaUsedV1(mediaId))
    }
  }

  for _, mediaId := range cmd.AttachmentMediaIds {
    media, err := s.mediaClient.FindById(ctx, mediaId)
    if err != nil {
      continue
    }
    attachments = append(attachments, vob.Attachment{
      IsEmbedded: false,
      Filename:   media.Name,
      Data:       media.Data,
    })

    if media.Type == vob.MediaTypeOneTime {
      events = append(events, intev.NewOneTimeMediaUsedV1(mediaId))
    }
  }

  return attachments, events
}

func (s *sendMailHandler) sendMail(ctx context.Context, span trace.Span, cmd *SendMailCommand, mail *entity.Mail) {
  ctx = trace.ContextWithSpan(ctx, span)

  // get all attachment medias
  attachments, oneTimeEvents := s.getAttachments(ctx, cmd)

  err := s.mailClient.Send(ctx, mail, attachments)
  isSucceed := err == nil

  var ev types.Event
  if isSucceed {
    // Succeed
    ev, err = mail.Delivered()
  } else {
    // Failed
    spanUtil.RecordError(err, span)
    ev, err = mail.Failed()
  }
  if err != nil {
    spanUtil.RecordError(err, span)
    return
  }

  // Update status
  err = mail.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return
  }

  // Update
  err = s.persistent.Update(ctx, mail)
  if err != nil {
    spanUtil.RecordError(err, span)
    return
  }

  // Add integration event
  mail.AddEvents(intev.NewMailDeliveredV1(mail.Id, isSucceed, time.Now()))
  mail.AddEvents(oneTimeEvents...)

  // Publish
  err = s.publisher.Publish(ctx, mail)
  if err != nil {
    spanUtil.RecordError(err, span)
  }
}

func (s *sendMailHandler) Handle(ctx context.Context, cmd *SendMailCommand) (types.Id, status.Object) {
  ctx, span := s.tracer.Start(ctx, "sendMailHandler.Handle")
  defer span.End()

  domain, err := cmd.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  mail, event := entity.CreateMail(&domain)

  err = mail.ApplyEvent(event)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrBadRequest(err)
  }

  err = s.persistent.Create(ctx, &mail)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository(err)
  }

  // Publish event
  err = s.publisher.Publish(ctx, &mail)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }
  mail.Clear()

  go s.sendMail(ctx, span, cmd, &mail)

  return mail.Id, status.Succeed()
}
