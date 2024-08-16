package mapper

import (
  "github.com/arcorium/rashop/proto/gen/go/common"
  mailerv1 "github.com/arcorium/rashop/proto/gen/go/mailer/v1"
  "github.com/arcorium/rashop/services/mailer/internal/app/dto"
  "github.com/arcorium/rashop/services/mailer/internal/app/query"
  sharedDto "github.com/arcorium/rashop/shared/dto"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  "google.golang.org/protobuf/types/known/timestamppb"
)

func ToGetMailsQuery(request *mailerv1.GetMailsRequest) query.GetMailsQuery {
  if request == nil {
    return query.GetMailsQuery{}
  }
  if request.Detail == nil {
    return query.GetMailsQuery{}
  }
  return query.GetMailsQuery{
    PagedElementDTO: sharedDto.PagedElementDTO{
      Element: request.Detail.Element,
      Page:    request.Detail.Page,
    },
  }
}

func ToFindMailByIdsQuery(request *mailerv1.FindMailByIdsRequest) (query.FindMailByIdsQuery, error) {
  mailIds, ierr := sharedUtil.CastSliceErrs(request.MailIds, types.IdFromString)
  if ierr.IsError() {
    return query.FindMailByIdsQuery{}, ierr.ToGRPCError("mail_ids")
  }

  return query.FindMailByIdsQuery{
    Ids: mailIds,
  }, nil
}

func ToProtoPagedElement[T any](result *sharedDto.PagedElementResult[T]) *common.PagedElementResponseDetail {
  return &common.PagedElementResponseDetail{
    Element:       result.Element,
    Page:          result.Page,
    TotalElements: result.TotalElements,
    TotalPages:    result.TotalPages,
  }
}

func ToProtoMail(mail *dto.MailResponseDTO) *mailerv1.Mail {
  var deliveredAt *timestamppb.Timestamp
  if !mail.DeliveredAt.IsZero() {
    deliveredAt = timestamppb.New(mail.DeliveredAt)
  }

  return &mailerv1.Mail{
    Id:          mail.Id.String(),
    Recipient:   sharedUtil.CastSlice(mail.Recipients, sharedUtil.ToString[types.Email]),
    Sender:      mail.Sender.String(),
    Tag:         toProtoTag(mail.Tag),
    Subject:     mail.Subject,
    Status:      toProtoStatus(mail.Status),
    SentAt:      timestamppb.New(mail.SentAt),
    DeliveredAt: deliveredAt,
  }
}
