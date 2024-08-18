package tracer

import (
  "github.com/arcorium/rashop/services/token/constant"
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/trace"
)

func Get(options ...trace.TracerOption) trace.Tracer {
  return otel.Tracer(constant.ServiceName, options...)
}
