package command

import (
  "github.com/arcorium/rashop/services/media_storage/internal/domain/repository"
  "github.com/arcorium/rashop/services/media_storage/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

type CommonHandlerParameter struct {
  Persistent repository.IMediaPersistent
  Storage    repository.IStorageClient
  Publisher  repository.IMessagePublisher
}

func newCommonHandler(parameter *CommonHandlerParameter, opts ...trace.TracerOption) commonHandler {
  return commonHandler{
    persistent: parameter.Persistent,
    storage:    parameter.Storage,
    publisher:  parameter.Publisher,
    tracer:     tracer.Get(opts...),
  }
}

type commonHandler struct {
  persistent repository.IMediaPersistent
  storage    repository.IStorageClient
  publisher  repository.IMessagePublisher
  tracer     trace.Tracer
}
