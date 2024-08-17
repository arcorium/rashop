package command

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/repository"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IDeleteMediaHandler interface {
  handler.CommandUnit[*DeleteMediaCommand]
}

func NewDeleteMediaHandler(parameter CommonHandlerParameter) IDeleteMediaHandler {
  return &deleteMediaHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type deleteMediaHandler struct {
  commonHandler
}

func (d *deleteMediaHandler) Handle(ctx context.Context, cmd *DeleteMediaCommand) status.Object {
  ctx, span := d.tracer.Start(ctx, "deleteMediaHandler.handle")
  defer span.End()

  media, err := d.persistent.FindByIds(ctx, cmd.MediaIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  var events []types.Event
  err = d.persistent.AsUnit(ctx, func(ctx context.Context, persistent repository.IMediaPersistent) error {
    err := persistent.Delete(ctx, cmd.MediaIds...)
    if err != nil {
      return err
    }

    for _, m := range media {
      // Event
      events = append(events, intev.NewMediaDeletedV1(m.Id, m.Usage.Underlying()))
      // Delete from storage
      err = d.storage.Delete(ctx, m.ProviderPath)
      if err != nil {
        return err
      }
    }
    return nil
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  err = d.publisher.PublishEvents(ctx, events...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Deleted()
}
