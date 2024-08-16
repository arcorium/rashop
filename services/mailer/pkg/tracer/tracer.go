package tracer

import (
  "github.com/arcorium/rashop/services/mailer/constant"
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/trace"
)

func Get(options ...trace.TracerOption) trace.Tracer {
  return otel.Tracer(constant.SERVICE_NAME, options...)
}
