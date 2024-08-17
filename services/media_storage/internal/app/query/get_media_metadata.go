package query

import (
  "context"
  "github.com/arcorium/rashop/services/media_storage/internal/app/dto"
  "github.com/arcorium/rashop/services/media_storage/internal/app/mapper"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IGetMediaMetadataHandler interface {
  handler.Command[*GetMediaMetadataQuery, []dto.MediaMetadataResponseDTO]
}

func NewGetMediaMetadataHandler(parameter CommonHandlerParameter) IGetMediaMetadataHandler {
  return &getMediaMetadataHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type getMediaMetadataHandler struct {
  commonHandler
}

func (g *getMediaMetadataHandler) Handle(ctx context.Context, query *GetMediaMetadataQuery) ([]dto.MediaMetadataResponseDTO, status.Object) {
  ctx, span := g.tracer.Start(ctx, "getMediaMetadataHandler.Handle")
  defer span.End()

  // Get metadata
  result, err := g.persistent.FindByIds(ctx, query.MediaIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err)
  }

  for i := range len(result) {
    media := &result[i]
    if len(media.FullPath) > 0 {
      continue
    }

    path, err := g.storage.GetFullPath(ctx, media.ProviderPath, media.IsPublic)
    if err != nil {
      spanUtil.RecordError(err, span)
      return nil, status.FromRepository(err)
    }
    media.FullPath = path
  }

  resp := sharedUtil.CastSliceP(result, mapper.ToMediaMetadataResponseDTO)
  return resp, status.Succeed()
}
