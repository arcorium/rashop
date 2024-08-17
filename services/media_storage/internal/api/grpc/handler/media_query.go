package handler

import (
  "context"
  "github.com/arcorium/rashop/proto/gen/go/media_storage/v1"
  "github.com/arcorium/rashop/services/media_storage/internal/api/grpc/mapper"
  "github.com/arcorium/rashop/services/media_storage/internal/app/service"
  "github.com/arcorium/rashop/services/media_storage/pkg/tracer"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
)

func NewMediaQuery(svc service.IMediaQuery) MediaQueryHandler {
  return MediaQueryHandler{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type MediaQueryHandler struct {
  mediav1.UnimplementedMediaStorageQueryServiceServer

  svc    service.IMediaQuery
  tracer trace.Tracer
}

func (m *MediaQueryHandler) Register(server *grpc.Server) {
  mediav1.RegisterMediaStorageQueryServiceServer(server, m)
}

func (m *MediaQueryHandler) Get(request *mediav1.GetMediaRequest, server mediav1.MediaStorageQueryService_GetServer) error {
  ctx, span := m.tracer.Start(server.Context(), "MediaQueryHandler.Get")
  defer span.End()

  query, err := mapper.ToGetMediaQuery(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  result, stat := m.svc.Get(ctx, &query)
  if stat.IsError() {
    spanUtil.RecordError(err, span)
    return stat.ToGRPCError()
  }

  // Send data back to client by streaming
  for {
    err = server.Send(&mediav1.GetMediaResponse{
      Name:       result.Name,
      MediaChunk: result.Data,
    })
    if err != nil {
      spanUtil.RecordError(err, span)
      return err
    }
    break
  }
  return nil
}

func (m *MediaQueryHandler) GetMetadata(ctx context.Context, request *mediav1.GetMediaMetadataRequest) (*mediav1.GetMediaMetadataResponse, error) {
  ctx, span := m.tracer.Start(ctx, "MediaQueryHandler.GetMetadata")
  defer span.End()

  query, err := mapper.ToGetMetadataQuery(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  metadata, stat := m.svc.GetMetadata(ctx, &query)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := mapper.ToProtoGetMediaMetadataResponse(metadata...)
  return resp, nil
}
