package query

import (
  "context"
  "github.com/arcorium/rashop/services/media_storage/internal/app/dto"
  "github.com/arcorium/rashop/services/media_storage/internal/app/mapper"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IGetMediaHandler interface {
  handler.Command[*GetMediaQuery, dto.MediaResponseDTO]
}

func NewGetMediaHandler(parameter CommonHandlerParameter) IGetMediaHandler {
  return &getMediaHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type getMediaHandler struct {
  commonHandler
}

func (g *getMediaHandler) Handle(ctx context.Context, query *GetMediaQuery) (dto.MediaResponseDTO, status.Object) {
  ctx, span := g.tracer.Start(ctx, "getMediaHandler.Handle")
  defer span.End()

  // Get metadata
  result, err := g.persistent.FindByIds(ctx, query.MediaId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.MediaResponseDTO{}, status.FromRepository(err)
  }

  // Get file
  current := &result[0]
  current, err = g.storage.Get(ctx, current)

  resp := mapper.ToMediaResponseDTO(current)
  return resp, status.Succeed()
}
