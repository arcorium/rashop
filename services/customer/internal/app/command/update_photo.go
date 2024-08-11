package command

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "mini-shop/services/user/pkg/cqrs"
)

type IUpdatePhotoHandler interface {
  handler.CommandUnit[*UpdateCustomerPhotoCommand]
}

func NewUpdatePhotoHandler(parameter cqrs.CommonHandlerParameter) IUpdatePhotoHandler {
  return &updatePhotoHandler{
    basicHandler: newBasicHandler(&parameter),
  }
}

type updatePhotoHandler struct {
  basicHandler
}

func (u *updatePhotoHandler) Handle(ctx context.Context, cmd *UpdateCustomerPhotoCommand) status.Object {
  ctx, span := u.tracer.Start(ctx, "updatePhotoHandler.Handle")
  defer span.End()

  customers, err := u.repo.FindByIds(ctx, cmd.CustomerId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  current := &customers[0]
  ev := current.UpdatePhoto(cmd.MediaId)

  current.AddEvents(
    intev.NewCustomerPhotoChangedV1(current.Id, current.PhotoId, cmd.MediaId),
  )

  err = current.ApplyEvent(ev)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err = u.repo.Update(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err)
  }

  err = u.publisher.PublishAggregate(ctx, current)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  return status.Updated()
}
