package query

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/app/dto"
  "github.com/arcorium/rashop/services/mailer/internal/app/mapper"
  sharedDto "github.com/arcorium/rashop/shared/dto"
  "github.com/arcorium/rashop/shared/interfaces/handler"
  "github.com/arcorium/rashop/shared/status"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
)

type IGetMailsHandler interface {
  handler.Command[*GetMailsQuery, sharedDto.PagedElementResult[dto.MailResponseDTO]]
}

func NewGetMailsHandler(parameter CommonHandlerParameter) IGetMailsHandler {
  return &getMailsHandler{
    commonHandler: newCommonHandler(&parameter),
  }
}

type getMailsHandler struct {
  commonHandler
}

func (g *getMailsHandler) Handle(ctx context.Context, cmd *GetMailsQuery) (sharedDto.PagedElementResult[dto.MailResponseDTO], status.Object) {
  ctx, span := g.tracer.Start(ctx, "getMailsHandler.Handle")
  defer span.End()

  result, err := g.persistent.Get(ctx, cmd.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.MailResponseDTO]{}, status.FromRepository(err)
  }

  resp := sharedUtil.CastSliceP(result.Data, mapper.ToMailResponse)
  return sharedDto.NewPagedElementResult2(resp, &cmd.PagedElementDTO, result.Total), status.Succeed()
}
