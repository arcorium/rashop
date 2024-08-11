package command

import (
  "go.opentelemetry.io/otel/trace"
  "mini-shop/services/user/internal/domain/repository"
  "mini-shop/services/user/pkg/cqrs"
  "mini-shop/services/user/pkg/tracer"
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
