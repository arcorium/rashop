package command

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/domain/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IDeleteMailsHandler interface {
  handler.Command[*DeleteMailsCommand, uint64]
}

func NewDeleteMailsHandler(parameter CommonHandlerParameter) IDeleteMailsHandler {
  return &deleteMailHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type deleteMailHandler struct {
  commonHandler
}

func (d *deleteMailHandler) Handle(ctx context.Context, cmd *DeleteMailsCommand) (uint64, status.Object) {
  ctx, span := d.tracer.Start(ctx, "deleteMailHandler.Handle")
  defer span.End()

  ev := &event.MailDeletedV1{
    DomainV1:  event.NewV1(),
    StartTime: cmd.StartTime,
    EndTime:   cmd.EndTime,
  }

  total, err := d.persistent.Delete(ctx, cmd.StartTime, cmd.EndTime)
  if err != nil {
    spanUtil.RecordError(err, span)
    return 0, status.FromRepository(err)
  }

  err = d.publisher.PublishEvents(ctx, ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return 0, status.ErrInternal(err)
  }

  return total, status.Deleted()
}
