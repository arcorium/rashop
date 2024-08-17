package handler

import (
  "context"
  "github.com/arcorium/rashop/proto/gen/go/media_storage/v1"
  "github.com/arcorium/rashop/services/media_storage/internal/api/grpc/mapper"
  "github.com/arcorium/rashop/services/media_storage/internal/app/command"
  "github.com/arcorium/rashop/services/media_storage/internal/app/service"
  "github.com/arcorium/rashop/services/media_storage/pkg/tracer"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "io"
)

func NewMediaCommand(svc service.IMediaCommand) MediaCommandHandler {
  return MediaCommandHandler{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type MediaCommandHandler struct {
  mediav1.UnimplementedMediaStorageCommandServiceServer

  svc    service.IMediaCommand
  tracer trace.Tracer
}

func (m *MediaCommandHandler) Register(server *grpc.Server) {
  mediav1.RegisterMediaStorageCommandServiceServer(server, m)
}

func (m *MediaCommandHandler) Store(server mediav1.MediaStorageCommandService_StoreServer) error {
  ctx, span := m.tracer.Start(server.Context(), "MediaCommandHandler.Store")
  defer span.End()

  cmd := command.StoreMediaCommand{}
  for {
    req, err := server.Recv()
    if err != nil {
      if err == io.EOF {
        break
      }
      spanUtil.RecordError(err, span)
      return err
    }

    cmd.Name = req.Name
    cmd.IsPublic = req.IsPublic
    cmd.Usage = mapper.ToDomainUsage(req.Usage)
    cmd.Data = append(cmd.Data, req.MediaChunk...)
  }

  id, stat := m.svc.Store(ctx, &cmd)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return stat.ToGRPCError()
  }

  resp := &mediav1.StoreMediaResponse{
    MediaId: id.String(),
  }
  return server.SendAndClose(resp)
}

func (m *MediaCommandHandler) Delete(ctx context.Context, request *mediav1.DeleteMediaRequest) (*mediav1.DeleteMediaResponse, error) {
  ctx, span := m.tracer.Start(ctx, "MediaCommandHandler.Delete")
  defer span.End()

  cmd, err := mapper.ToDeleteMediaCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := m.svc.Delete(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
