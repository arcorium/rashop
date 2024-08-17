package consumer

import (
  "context"
  intev "github.com/arcorium/rashop/contract/integration/event"
  "github.com/arcorium/rashop/services/media_storage/internal/app/service"
  "github.com/arcorium/rashop/services/media_storage/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

func NewMediaCommandHandler(svc service.IMediaCommandConsumer) IMediaCommandHandler {
  return &mediaCommandConsumerHandler{
    svc:    svc,
    tracer: tracer.Get(),
  }
}

type IMediaCommandHandler interface {
  OnOneTimeMediaUsed(ctx context.Context, v1 *intev.OneTimeMediaUsedV1) error
}

type mediaCommandConsumerHandler struct {
  svc    service.IMediaCommandConsumer
  tracer trace.Tracer
}

func (m *mediaCommandConsumerHandler) OnOneTimeMediaUsed(ctx context.Context, ev *intev.OneTimeMediaUsedV1) error {
  ctx, span := m.tracer.Start(ctx, "mediaCommandConsumerHandler.OnOneTimeMediaUsed")
  defer span.End()

  stat := m.svc.OneTimeMediaUsed(ctx, ev)
  return stat.ErrorWithSpan(span)
}
