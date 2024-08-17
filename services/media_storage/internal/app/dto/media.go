package dto

import (
  vob "github.com/arcorium/rashop/services/media_storage/internal/domain/valueobject"
  "github.com/arcorium/rashop/shared/types"
  "time"
)

type MediaResponseDTO struct {
  Name string
  Data []byte
}

type MediaMetadataResponseDTO struct {
  Id             types.Id
  Name           string
  Usage          vob.MediaUsage
  ContentType    string
  Size           uint64
  Url            string
  IsPublic       bool
  LastModifiedAt time.Time
  StoredAt       time.Time
}
