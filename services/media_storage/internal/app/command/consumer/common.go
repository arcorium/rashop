package consumer

import (
  "github.com/arcorium/rashop/services/media_storage/internal/app/command"
  "github.com/arcorium/rashop/services/media_storage/internal/domain/repository"
  "github.com/arcorium/rashop/services/media_storage/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

type CommonHandlerParameter struct {
  Persistent repository.IMediaPersistent
  Delete     command.IDeleteMediaHandler
}

func newCommonHandler(parameter *CommonHandlerParameter, opts ...trace.TracerOption) commonHandler {
  return commonHandler{
    persistent: parameter.Persistent,
    delete:     parameter.Delete,
    tracer:     tracer.Get(opts...),
  }
}

type commonHandler struct {
  persistent repository.IMediaPersistent
  delete     command.IDeleteMediaHandler
  tracer     trace.Tracer
}
