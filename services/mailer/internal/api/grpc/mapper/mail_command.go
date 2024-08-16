package mapper

import (
  "fmt"
  mailerv1 "github.com/arcorium/rashop/proto/gen/go/mailer/v1"
  "github.com/arcorium/rashop/services/mailer/internal/app/command"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
  sharedErr "github.com/arcorium/rashop/shared/errors"
  "github.com/arcorium/rashop/shared/status"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
)

func ToSendMailCommand(request *mailerv1.SendMailRequest) (command.SendMailCommand, error) {
  var errors []sharedErr.FieldError

  tag, err := vob.NewMailTag(uint8(request.Tag))
  if err != nil {
    errors = append(errors, sharedErr.NewFieldError("tag", err))
  }

  emails, ierr := sharedUtil.CastSliceErrs(request.Recipients, types.EmailFromString)
  if ierr.IsError() {
    errors = append(errors, sharedErr.NewFieldError("recipients", ierr))
  }

  var sender *types.Email
  if request.Sender != nil {
    temp, err := types.EmailFromString(*request.Sender)
    if err != nil {
      errors = append(errors, sharedErr.NewFieldError("sender", err))
    } else {
      sender = &temp
    }
  }

  bodyType, err := vob.NewBodyType(uint8(request.BodyType))
  if err != nil {
    errors = append(errors, sharedErr.NewFieldError("body_type", err))
  }

  embedIds, ierr := sharedUtil.CastSliceErrs(request.EmbeddedMediaIds, types.IdFromString)
  if ierr.IsError() && !ierr.IsEmptySlice() {
    errors = append(errors, sharedErr.NewFieldError("embedded_media_ids", ierr))
  }

  attachmentIds, ierr := sharedUtil.CastSliceErrs(request.AttachmentMediaIds, types.IdFromString)
  if ierr.IsError() && !ierr.IsEmptySlice() {
    errors = append(errors, sharedErr.NewFieldError("attachment_media_ids", ierr))
  }

  if len(errors) > 0 {
    return command.SendMailCommand{}, sharedErr.GrpcFieldErrors(errors...)
  }

  return command.SendMailCommand{
    Tag:                tag,
    Recipients:         emails,
    Sender:             types.NewNullable(sender),
    Subject:            request.Subject,
    BodyType:           bodyType,
    Body:               request.Body,
    EmbeddedMediaIds:   embedIds,
    AttachmentMediaIds: attachmentIds,
  }, nil
}

func ToDeleteMailCommand(request *mailerv1.DeleteMailRequest) (command.DeleteMailsCommand, error) {
  result := command.DeleteMailsCommand{}

  flag := false
  if request.StartTime == nil {
    flag = true
  } else {
    result.StartTime = request.StartTime.AsTime()
  }

  if request.UntilTime == nil && flag {
    stat := status.ErrBadRequest(fmt.Errorf("arguments can't be nulled for both, it should be present at least one"))
    return command.DeleteMailsCommand{}, stat.ToGRPCError()
  } else {
    result.EndTime = request.UntilTime.AsTime()
  }

  return result, nil
}
