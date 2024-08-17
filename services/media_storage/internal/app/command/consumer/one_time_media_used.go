package consumer

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/media_storage/internal/app/mapper"
  vob "github.com/arcorium/rashop/services/media_storage/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IOneTimeMediaUsedHandler interface {
  handler.Consumer[*intev.OneTimeMediaUsedV1]
}

func NewOneTimeMediaUsedHandler(parameter CommonHandlerParameter) IOneTimeMediaUsedHandler {
  return &oneTimeMediaUsedHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type oneTimeMediaUsedHandler struct {
  commonHandler
}

func (o *oneTimeMediaUsedHandler) Handle(ctx context.Context, ev *intev.OneTimeMediaUsedV1) status.Object {
  ctx, span := o.tracer.Start(ctx, "oneTimeMediaUsedHandler.Handle")
  defer span.End()

  cmd, err := mapper.MapOneTimeMediaUsedEventToCommand(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  media, err := o.persistent.FindByIds(ctx, cmd.MediaIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &media[0]
  if current.Usage != vob.UsageOnce {
    return status.Succeed()
  }

  return o.delete.Handle(ctx, &cmd)
}
