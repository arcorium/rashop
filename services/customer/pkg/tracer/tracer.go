package tracer

import (
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/trace"
  "mini-shop/services/user/constant"
)

func Get(options ...trace.TracerOption) trace.Tracer {
  return otel.Tracer(constant.SERVICE_NAME, options...)
}
