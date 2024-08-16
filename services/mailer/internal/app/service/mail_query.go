package service

import (
  "context"
  "github.com/arcorium/rashop/services/mailer/internal/app/dto"
  "github.com/arcorium/rashop/services/mailer/internal/app/query"
  sharedDto "github.com/arcorium/rashop/shared/dto"
  "github.com/arcorium/rashop/shared/status"
)

type IMailQuery interface {
  Get(ctx context.Context, query *query.GetMailsQuery) (sharedDto.PagedElementResult[dto.MailResponseDTO], status.Object)
  FindByIds(ctx context.Context, query *query.FindMailByIdsQuery) ([]dto.MailResponseDTO, status.Object)
}

func NewMailQuery(config MailQueryConfig) IMailQuery {
  return &mailQueryService{
    i: config,
  }
}

func DefaultMailQueryConfig(parameter query.CommonHandlerParameter) MailQueryConfig {
  return MailQueryConfig{
    Get:       query.NewGetMailsHandler(parameter),
    FindByIds: query.NewFindMailByIdsHandler(parameter),
  }
}

type MailQueryConfig struct {
  Get       query.IGetMailsHandler
  FindByIds query.IFindMailByIdsHandler
}

type mailQueryService struct {
  i MailQueryConfig
}

func (m *mailQueryService) Get(ctx context.Context, query *query.GetMailsQuery) (sharedDto.PagedElementResult[dto.MailResponseDTO], status.Object) {
  return m.i.Get.Handle(ctx, query)
}

func (m *mailQueryService) FindByIds(ctx context.Context, query *query.FindMailByIdsQuery) ([]dto.MailResponseDTO, status.Object) {
  return m.i.FindByIds.Handle(ctx, query)
}
