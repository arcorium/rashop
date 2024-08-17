package mapper

import (
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/media_storage/internal/app/command"
  "github.com/arcorium/rashop/services/media_storage/internal/app/dto"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/entity"
  "github.com/arcorium/rashop/shared/types"
)

func ToMediaResponseDTO(media *entity.Media) dto.MediaResponseDTO {
  return dto.MediaResponseDTO{
    Name: media.Name,
    Data: media.Data,
  }
}

func ToMediaMetadataResponseDTO(media *entity.Media) dto.MediaMetadataResponseDTO {
  return dto.MediaMetadataResponseDTO{
    Id:             media.Id,
    Name:           media.Name,
    Usage:          media.Usage,
    ContentType:    media.ContentType,
    Size:           media.Size,
    Url:            media.FullPath,
    IsPublic:       media.IsPublic,
    LastModifiedAt: media.LastModifiedAt,
    StoredAt:       media.StoredAt,
  }
}

func MapOneTimeMediaUsedEventToCommand(ev *intev.OneTimeMediaUsedV1) (command.DeleteMediaCommand, error) {
  id, err := types.IdFromString(ev.MediaId)
  if err != nil {
    return command.DeleteMediaCommand{}, err
  }

  return command.DeleteMediaCommand{
    MediaIds: []types.Id{id},
  }, nil
}
