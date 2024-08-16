package consumer

import (
  "github.com/arcorium/rashop/services/mailer/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

func newCommonHandler(opts ...trace.TracerOption) commonHandler {
  return commonHandler{
    tracer: tracer.Get(opts...),
  }
}

type commonHandler struct {
  tracer trace.Tracer
}
