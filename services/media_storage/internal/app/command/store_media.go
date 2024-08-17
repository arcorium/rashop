package command

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/repository"
  "github.com/arcorium/rashop/shared/grpc"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IStoreMediaHandler interface {
  handler.Command[*StoreMediaCommand, types.Id]
}

func NewStoreMediaHandler(parameter CommonHandlerParameter) IStoreMediaHandler {
  return &storeMediaHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type storeMediaHandler struct {
  commonHandler
}

func (g *storeMediaHandler) Handle(ctx context.Context, cmd *StoreMediaCommand) (types.Id, status.Object) {
  ctx, span := g.tracer.Start(ctx, "storeMediaHandler.Handle")
  defer span.End()

  userId, err := grpc.ExtractUserId(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  // Create metadata
  media, err := cmd.ToDomain(g.storage)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrBadRequest(err)
  }

  err = g.persistent.AsUnit(ctx, func(ctx context.Context, persistent repository.IMediaPersistent) error {
    defer func() {
      // Handle media deletion when error occurred
      if err != nil {
        err = g.storage.Delete(ctx, media.ProviderPath)
      }
    }()

    // Upload media
    err := g.storage.Store(ctx, &media)
    if err != nil {
      return err
    }

    err = persistent.Create(ctx, &media)
    return err
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository(err)
  }

  event := intev.NewMediaStoredV1(userId, media.Id, media.Usage.Underlying())
  err = g.publisher.PublishEvents(ctx, event)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  return media.Id, status.Created()
}
