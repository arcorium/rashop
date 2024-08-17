package query

import (
  "github.com/arcorium/rashop/services/media_storage/internal/domain/repository"
  "github.com/arcorium/rashop/services/media_storage/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

type CommonHandlerParameter struct {
  Persistent repository.IMediaPersistent
  Storage    repository.IStorageClient
}

func newCommonHandler(parameter *CommonHandlerParameter, opts ...trace.TracerOption) commonHandler {
  return commonHandler{
    persistent: parameter.Persistent,
    storage:    parameter.Storage,
    tracer:     tracer.Get(opts...),
  }
}

type commonHandler struct {
  persistent repository.IMediaPersistent
  storage    repository.IStorageClient
  tracer     trace.Tracer
}
