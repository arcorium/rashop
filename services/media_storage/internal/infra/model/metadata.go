package model

import (
  "github.com/arcorium/rashop/services/media_storage/internal/domain/entity"
  vob "github.com/arcorium/rashop/services/media_storage/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  "github.com/arcorium/rashop/shared/util/repo"
  "github.com/uptrace/bun"
  "time"
)

type MetadataMapOption = repo.DataAccessModelMapOption[*entity.Media, *Metadata]

func FromMediaDomain(media *entity.Media, opts ...MetadataMapOption) Metadata {
  md := Metadata{
    Id:           media.Id.String(),
    Name:         media.Name,
    MimeType:     media.ContentType,
    Size:         media.Size,
    IsPublic:     media.IsPublic,
    Provider:     media.Provider.Underlying(),
    ProviderPath: media.ProviderPath,
    UpdatedAt:    media.LastModifiedAt,
    CreatedAt:    media.StoredAt,
  }

  for _, opt := range opts {
    opt(media, &md)
  }
  return md
}

type Metadata struct {
  bun.BaseModel `bun:"table:media_metadata,alias:mm"`
  Id            string    `bun:",nullzero,type:uuid,pk"`
  Name          string    `bun:",nullzero,notnull"`
  MimeType      string    `bun:",nullzero"`
  Size          uint64    `bun:",notnull"`
  IsPublic      bool      `bun:",notnull"`
  Provider      uint8     `bun:",notnull"`
  ProviderPath  string    `bun:",nullzero,notnull"`
  UpdatedAt     time.Time `bun:",nullzero"`
  CreatedAt     time.Time `bun:",nullzero,notnull"`
}

func (m *Metadata) ToDomain() (entity.Media, error) {
  id, err := types.IdFromString(m.Id)
  if err != nil {
    return entity.Media{}, err
  }

  provider, err := vob.NewStorageProvider(m.Provider)
  if err != nil {
    return entity.Media{}, err
  }

  return entity.Media{
    Id:             id,
    Name:           m.Name,
    ContentType:    m.MimeType,
    Size:           m.Size,
    IsPublic:       m.IsPublic,
    Provider:       provider,
    ProviderPath:   m.ProviderPath,
    StoredAt:       m.CreatedAt,
    LastModifiedAt: m.UpdatedAt,
  }, nil
}
