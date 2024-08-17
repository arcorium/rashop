package command

import (
  "github.com/arcorium/rashop/services/media_storage/internal/domain/entity"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/repository"
  vob "github.com/arcorium/rashop/services/media_storage/internal/domain/valueobject"
  "github.com/arcorium/rashop/services/media_storage/pkg/util"
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type StoreMediaCommand struct {
  Name     string
  IsPublic bool
  Usage    vob.MediaUsage
  Data     []byte
}

func (s *StoreMediaCommand) ToDomain(provider repository.IStorageClient) (entity.Media, error) {
  id, err := types.NewId()
  if err != nil {
    return entity.Media{}, err
  }

  return entity.Media{
    Id:          id,
    Name:        s.Name,
    Usage:       s.Usage,
    Data:        s.Data,
    ContentType: util.GetMimeType(s.Name),
    Size:        uint64(len(s.Data)),
    IsPublic:    s.IsPublic,
    Provider:    provider.GetProvider(),
    StoredAt:    time.Now(),
  }, nil
}

type DeleteMediaCommand struct {
  MediaIds []types.Id
}
