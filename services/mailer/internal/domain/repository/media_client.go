package repository

import (
  "context"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
)

type IMediaStorageClient interface {
  FindById(ctx context.Context, mediaId types.Id) (vob.Media, error)
}
