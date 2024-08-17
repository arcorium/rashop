package repository

import (
  "context"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/entity"
  "github.com/arcorium/rashop/shared/types"
  "github.com/arcorium/rashop/shared/util/repo"
)

type IMediaPersistent interface {
  repo.IUnitOfWork[IMediaPersistent]
  Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Media], error)
  FindByIds(ctx context.Context, mediaIds ...types.Id) ([]entity.Media, error)
  Create(ctx context.Context, media *entity.Media) error
  Update(ctx context.Context, media *entity.Media) error
  Delete(ctx context.Context, mediaIds ...types.Id) error
}
