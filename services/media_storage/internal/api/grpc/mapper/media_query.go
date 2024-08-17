package mapper

import (
  mediav1 "github.com/arcorium/rashop/proto/gen/go/media_storage/v1"
  "github.com/arcorium/rashop/services/media_storage/internal/app/dto"
  "github.com/arcorium/rashop/services/media_storage/internal/app/query"
  sharedErr "github.com/arcorium/rashop/shared/errors"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "google.golang.org/protobuf/types/known/timestamppb"
)

func ToGetMediaQuery(request *mediav1.GetMediaRequest) (query.GetMediaQuery, error) {
  mediaId, err := types.IdFromString(request.MediaId)
  if err != nil {
    return query.GetMediaQuery{}, sharedErr.NewFieldError("media_id", err).ToGrpcError()
  }

  return query.GetMediaQuery{
    MediaId: mediaId,
  }, nil
}

func ToGetMetadataQuery(request *mediav1.GetMediaMetadataRequest) (query.GetMediaMetadataQuery, error) {
  ids, ierr := sharedUtil.CastSliceErrs(request.MediaIds, types.IdFromString)
  if ierr.IsError() {
    return query.GetMediaMetadataQuery{}, sharedErr.NewFieldError("media_ids", ierr).ToGrpcError()
  }

  return query.GetMediaMetadataQuery{
    MediaIds: ids,
  }, nil
}

func ToProtoMediaMetadata(dto *dto.MediaMetadataResponseDTO) *mediav1.MediaMetadata {
  var lastModifiedAt *timestamppb.Timestamp
  if !dto.LastModifiedAt.IsZero() {
    lastModifiedAt = timestamppb.New(dto.LastModifiedAt)
  }

  return &mediav1.MediaMetadata{
    Id:           dto.Id.String(),
    Name:         dto.Name,
    Usage:        toProtoUsage(dto.Usage),
    ContentType:  dto.ContentType,
    Size:         dto.Size,
    Url:          dto.Url,
    IsPublic:     dto.IsPublic,
    LastModified: lastModifiedAt,
    StoredAt:     timestamppb.New(dto.StoredAt),
  }
}

func ToProtoGetMediaMetadataResponse(dtos ...dto.MediaMetadataResponseDTO) *mediav1.GetMediaMetadataResponse {
  result := make(map[string]*mediav1.MediaMetadata)
  for _, dto := range dtos {
    result[dto.Id.String()] = ToProtoMediaMetadata(&dto)
  }
  return &mediav1.GetMediaMetadataResponse{
    Medias: result,
  }
}
