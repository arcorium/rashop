package handler

import (
  "context"
  "github.com/arcorium/rashop/proto/gen/go/mailer/v1"
  "github.com/arcorium/rashop/services/mailer/internal/api/grpc/mapper"
  "github.com/arcorium/rashop/services/mailer/internal/app/service"
  "github.com/arcorium/rashop/services/mailer/pkg/tracer"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
)

func NewMailCommand(svc service.IMailCommand) MailCommand {
  return MailCommand{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type MailCommand struct {
  mailerv1.UnimplementedMailerCommandServiceServer

  svc    service.IMailCommand
  tracer trace.Tracer
}

func (m *MailCommand) Register(server *grpc.Server) {
  mailerv1.RegisterMailerCommandServiceServer(server, m)
}

func (m *MailCommand) Send(ctx context.Context, request *mailerv1.SendMailRequest) (*mailerv1.SendMailResponse, error) {
  ctx, span := m.tracer.Start(ctx, "MailCommand.Send")
  defer span.End()

  cmd, err := mapper.ToSendMailCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  id, stat := m.svc.Send(ctx, &cmd)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.SendMailResponse{
    MailId: id.String(),
  }
  return resp, nil
}

func (m *MailCommand) Delete(ctx context.Context, request *mailerv1.DeleteMailRequest) (*mailerv1.DeleteMailResponse, error) {
  ctx, span := m.tracer.Start(ctx, "MailCommand.Delete")
  defer span.End()

  cmd, err := mapper.ToDeleteMailCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  total, stat := m.svc.Delete(ctx, &cmd)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.DeleteMailResponse{
    DeletedCount: total,
  }
  return resp, nil
}
