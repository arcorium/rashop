package command

import (
  "context"
  "github.com/arcorium/rashop/services/token/internal/domain/entity"
  "github.com/arcorium/rashop/services/token/internal/domain/repository"
  "github.com/arcorium/rashop/services/token/pkg/tracer"
  "github.com/arcorium/rashop/shared/status"
  spanUtil "github.com/arcorium/rashop/shared/util/span"
  "go.opentelemetry.io/otel/trace"
)

type CommonHandlerParameter struct {
  Persistent repository.ITokenPersistent
  Publisher  repository.IMessagePublisher
}

func newCommonHandler(parameter *CommonHandlerParameter, opts ...trace.TracerOption) commonHandler {
  return commonHandler{
    persistent: parameter.Persistent,
    publisher:  parameter.Publisher,
    tracer:     tracer.Get(opts...),
  }
}

type commonHandler struct {
  persistent repository.ITokenPersistent
  publisher  repository.IMessagePublisher
  tracer     trace.Tracer
}

func (c *commonHandler) Publish(ctx context.Context, aggregate *entity.Token, successStatus status.Object) status.Object {
  span := trace.SpanFromContext(ctx)

  err := c.publisher.Publish(ctx, aggregate)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }
  return successStatus
}
