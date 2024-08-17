package mapper

import (
  mediav1 "github.com/arcorium/rashop/proto/gen/go/media_storage/v1"
  "github.com/arcorium/rashop/services/media_storage/internal/app/command"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
)

func ToDeleteMediaCommand(request *mediav1.DeleteMediaRequest) (command.DeleteMediaCommand, error) {
  ids, ierr := sharedUtil.CastSliceErrs(request.MediaIds, types.IdFromString)
  if ierr.IsError() {
    return command.DeleteMediaCommand{}, ierr.ToGRPCError("media_ids")
  }

  return command.DeleteMediaCommand{
    MediaIds: ids,
  }, nil
}
