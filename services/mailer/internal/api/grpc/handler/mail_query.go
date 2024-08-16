package handler

import (
  "context"
  "github.com/arcorium/rashop/proto/gen/go/mailer/v1"
  "github.com/arcorium/rashop/services/mailer/internal/api/grpc/mapper"
  "github.com/arcorium/rashop/services/mailer/internal/app/service"
  "github.com/arcorium/rashop/services/mailer/pkg/tracer"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
)

func NewMailQuery(svc service.IMailQuery) MailQuery {
  return MailQuery{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type MailQuery struct {
  mailerv1.UnimplementedMailerQueryServiceServer

  svc    service.IMailQuery
  tracer trace.Tracer
}

func (m *MailQuery) Register(server *grpc.Server) {
  mailerv1.RegisterMailerQueryServiceServer(server, m)
}

func (m *MailQuery) Get(ctx context.Context, request *mailerv1.GetMailsRequest) (*mailerv1.GetMailsResponse, error) {
  ctx, span := m.tracer.Start(ctx, "MailQuery.Get")
  defer span.End()

  query := mapper.ToGetMailsQuery(request)
  result, stat := m.svc.Get(ctx, &query)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.GetMailsResponse{
    Mails:  sharedUtil.CastSliceP(result.Data, mapper.ToProtoMail),
    Detail: mapper.ToProtoPagedElement(&result),
  }
  return resp, nil
}

func (m *MailQuery) FindByIds(ctx context.Context, request *mailerv1.FindMailByIdsRequest) (*mailerv1.FindMailByIdsResponse, error) {
  ctx, span := m.tracer.Start(ctx, "MailQuery.FindByIds")
  defer span.End()

  query, err := mapper.ToFindMailByIdsQuery(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, stat := m.svc.FindByIds(ctx, &query)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.FindMailByIdsResponse{
    Mails: sharedUtil.CastSliceP(result, mapper.ToProtoMail),
  }
  return resp, nil
}
