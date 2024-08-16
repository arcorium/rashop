package tracer

import (
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/trace"
  "rashop/services/customer/constant"
)

func Get(options ...trace.TracerOption) trace.Tracer {
  return otel.Tracer(constant.SERVICE_NAME, options...)
}
