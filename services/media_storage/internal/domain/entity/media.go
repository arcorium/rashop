package entity

import (
  vob "github.com/arcorium/rashop/services/media_storage/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type Media struct {
  types.AggregateBase
  types.AggregateHelper

  Id             types.Id
  Name           string
  Usage          vob.MediaUsage
  Data           []byte
  ContentType    string
  Size           uint64
  IsPublic       bool
  Provider       vob.StorageProvider
  ProviderPath   string
  FullPath       string
  StoredAt       time.Time
  LastModifiedAt time.Time
}
