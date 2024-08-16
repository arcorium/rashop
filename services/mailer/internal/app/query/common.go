package query

import (
  "github.com/arcorium/rashop/services/mailer/internal/domain/repository"
  "github.com/arcorium/rashop/services/mailer/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

type CommonHandlerParameter struct {
  Persistent repository.IMailPersistent
}

func newCommonHandler(parameter *CommonHandlerParameter, opts ...trace.TracerOption) commonHandler {
  return commonHandler{
    persistent: parameter.Persistent,
    tracer:     tracer.Get(opts...),
  }
}

type commonHandler struct {
  persistent repository.IMailPersistent
  tracer     trace.Tracer
}
