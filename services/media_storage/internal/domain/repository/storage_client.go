package repository

import (
  "context"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/entity"
  vob "github.com/arcorium/rashop/services/media_storage/internal/domain/valueobject"
)

type IStorageClient interface {
  Get(ctx context.Context, media *entity.Media) (*entity.Media, error)
  Store(ctx context.Context, media *entity.Media) error
  Delete(ctx context.Context, name string) error
  GetFullPath(ctx context.Context, name string, public bool) (string, error)
  GetProvider() vob.StorageProvider
}
