package consumer

import (
  "github.com/arcorium/rashop/services/token/internal/app/command"
  "github.com/arcorium/rashop/services/token/internal/domain/repository"
  "github.com/arcorium/rashop/services/token/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

type CommonHandlerParameter struct {
  Persistent repository.ITokenPersistent
  Generate   command.IGenerateTokenHandler
}

func newCommonHandler(parameter *CommonHandlerParameter, opts ...trace.TracerOption) commonHandler {
  return commonHandler{
    generate: parameter.Generate,
    tracer:   tracer.Get(opts...),
  }
}

type commonHandler struct {
  generate command.IGenerateTokenHandler
  tracer   trace.Tracer
}
