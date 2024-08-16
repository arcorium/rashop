package command

import (
  "github.com/arcorium/rashop/services/mailer/internal/domain/repository"
  "github.com/arcorium/rashop/services/mailer/pkg/tracer"
  "go.opentelemetry.io/otel/trace"
)

type CommonHandlerParameter struct {
  Persistent  repository.IMailPersistent
  MailClient  repository.IMailClient
  MediaClient repository.IMediaStorageClient
  Publisher   repository.IMessagePublisher
}

func newCommonHandler(parameter *CommonHandlerParameter, opts ...trace.TracerOption) commonHandler {
  return commonHandler{
    persistent:  parameter.Persistent,
    mailClient:  parameter.MailClient,
    mediaClient: parameter.MediaClient,
    publisher:   parameter.Publisher,
    tracer:      tracer.Get(opts...),
  }
}

type commonHandler struct {
  persistent  repository.IMailPersistent
  mailClient  repository.IMailClient
  mediaClient repository.IMediaStorageClient
  publisher   repository.IMessagePublisher
  tracer      trace.Tracer
}
