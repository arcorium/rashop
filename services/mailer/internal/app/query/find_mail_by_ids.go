package query

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/app/dto"
  "github.com/arcorium/rashop/services/mailer/internal/app/mapper"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IFindMailByIdsHandler interface {
  handler.Command[*FindMailByIdsQuery, []dto.MailResponseDTO]
}

func NewFindMailByIdsHandler(parameter CommonHandlerParameter) IFindMailByIdsHandler {
  return &findMailByIdsHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type findMailByIdsHandler struct {
  commonHandler
}

func (f *findMailByIdsHandler) Handle(ctx context.Context, cmd *FindMailByIdsQuery) ([]dto.MailResponseDTO, status.Object) {
  ctx, span := f.tracer.Start(ctx, "findMailByIdsHandler.Handle")
  defer span.End()

  result, err := f.persistent.FindByIds(ctx, cmd.Ids...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err)
  }

  resp := sharedUtil.CastSliceP(result, mapper.ToMailResponse)
  return resp, status.Succeed()
}
