package handler

import (
  "context"
  "github.com/arcorium/rashop/proto/gen/go/token/v1"
  "github.com/arcorium/rashop/services/token/internal/api/grpc/mapper"
  "github.com/arcorium/rashop/services/token/internal/app/service"
  "github.com/arcorium/rashop/services/token/pkg/tracer"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
)

func NewTokenCommand(svc service.ITokenCommand) TokenCommandHandler {
  return TokenCommandHandler{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type TokenCommandHandler struct {
  tokenv1.UnimplementedTokenCommandServiceServer

  svc    service.ITokenCommand
  tracer trace.Tracer
}

func (t *TokenCommandHandler) Register(server *grpc.Server) {
  tokenv1.RegisterTokenCommandServiceServer(server, t)
}

func (t *TokenCommandHandler) Generate(ctx context.Context, request *tokenv1.GenerateTokenRequest) (*tokenv1.GenerateTokenResponse, error) {
  ctx, span := t.tracer.Start(ctx, "TokenCommandHandler.Generate")
  defer span.End()

  cmd, err := mapper.ToGenerateTokenCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, stat := t.svc.Generate(ctx, &cmd)
  if stat.IsError() {
    spanUtil.RecordError(err, span)
    return nil, stat.Error
  }

  resp := &tokenv1.GenerateTokenResponse{
    Token: mapper.ToProtoToken(&result),
  }
  return resp, nil
}

func (t *TokenCommandHandler) Validate(ctx context.Context, request *tokenv1.ValidateTokenRequest) (*tokenv1.ValidateTokenResponse, error) {
  ctx, span := t.tracer.Start(ctx, "TokenCommandHandler.Validate")
  defer span.End()

  cmd, err := mapper.ToVerifyTokenCommand(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := t.svc.Verify(ctx, &cmd)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
