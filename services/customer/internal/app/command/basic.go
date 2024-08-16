package command

import (
  "go.opentelemetry.io/otel/trace"
  "rashop/services/customer/internal/domain/repository"
  "rashop/services/customer/pkg/cqrs"
  "rashop/services/customer/pkg/tracer"
)

func newBasicHandler(parameter *cqrs.CommonHandlerParameter, opts ...trace.TracerOption) basicHandler {
  return basicHandler{
    repo:      parameter.Repo,
    publisher: parameter.Publisher,
    tracer:    tracer.Get(opts...),
  }
}

type basicHandler struct {
  repo      repository.ICustomer
  publisher repository.IMessagePublisher
  tracer    trace.Tracer
}
