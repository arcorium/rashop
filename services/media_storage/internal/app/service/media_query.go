package service

import (
  "context"
  "github.com/arcorium/rashop/services/media_storage/internal/app/dto"
  "github.com/arcorium/rashop/services/media_storage/internal/app/query"
  "github.com/arcorium/rashop/shared/status"
)

type IMediaQuery interface {
  Get(ctx context.Context, query *query.GetMediaQuery) (dto.MediaResponseDTO, status.Object)
  GetMetadata(ctx context.Context, query *query.GetMediaMetadataQuery) ([]dto.MediaMetadataResponseDTO, status.Object)
}

func NewMediaQuery(config MediaQueryConfig) IMediaQuery {
  return &mediaQueryService{
    i: config,
  }
}

func DefaultMediaQueryConfig(parameter query.CommonHandlerParameter) MediaQueryConfig {
  return MediaQueryConfig{
    Get:         query.NewGetMediaHandler(parameter),
    GetMetadata: query.NewGetMediaMetadataHandler(parameter),
  }
}

type MediaQueryConfig struct {
  Get         query.IGetMediaHandler
  GetMetadata query.IGetMediaMetadataHandler
}

type mediaQueryService struct {
  i MediaQueryConfig
}

func (m *mediaQueryService) Get(ctx context.Context, query *query.GetMediaQuery) (dto.MediaResponseDTO, status.Object) {
  return m.i.Get.Handle(ctx, query)
}

func (m *mediaQueryService) GetMetadata(ctx context.Context, query *query.GetMediaMetadataQuery) ([]dto.MediaMetadataResponseDTO, status.Object) {
  return m.i.GetMetadata.Handle(ctx, query)
}
