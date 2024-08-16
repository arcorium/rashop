package otel

import (
  "context"
  "errors"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/shared/types"
  "github.com/dnwe/otelsarama"
  "github.com/prometheus/client_golang/prometheus"
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
  promExp "go.opentelemetry.io/otel/exporters/prometheus"
  "go.opentelemetry.io/otel/propagation"
  sdkmetric "go.opentelemetry.io/otel/sdk/metric"
  "go.opentelemetry.io/otel/sdk/resource"
  sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func newTraceProvider(ctx context.Context, resource *resource.Resource, opts ...otlptracegrpc.Option) (*sdktrace.TracerProvider, error) {
  exporter, err := otlptracegrpc.New(ctx, opts...)
  if err != nil {
    return nil, err
  }

  bsp := sdktrace.NewBatchSpanProcessor(exporter)
  traceProvider := sdktrace.NewTracerProvider(
    sdktrace.WithSampler(sdktrace.AlwaysSample()),
    sdktrace.WithResource(resource),
    sdktrace.WithSpanProcessor(bsp),
  )

  return traceProvider, nil
}

func newMetricProvider(resource *resource.Resource, collectors ...prometheus.Collector) (*sdkmetric.MeterProvider, error) {
  // Register collectors
  for _, collector := range collectors {
    err := prometheus.DefaultRegisterer.Register(collector)
    if err != nil {
      return nil, err
    }
  }

  // prometheus will use default registerer when not registering the registerer
  exporter, err := promExp.New()
  if err != nil {
    return nil, err
  }

  provider := sdkmetric.NewMeterProvider(
    sdkmetric.WithReader(exporter),
    sdkmetric.WithResource(resource),
  )

  return provider, nil
}

func newPropagator() propagation.TextMapPropagator {
  return propagation.NewCompositeTextMapPropagator(
    propagation.TraceContext{},
    propagation.Baggage{},
  )
}

type ShutdownFunc func(ctx context.Context) error

type SetupParameter struct {
  Resources  []attribute.KeyValue
  Options    []otlptracegrpc.Option
  Collectors []prometheus.Collector
}

// Setup trace, propagator, and metrics and set it as global
func Setup(parameter SetupParameter) (ShutdownFunc, error) {
  var err error
  var shutdowns []ShutdownFunc

  // Resource
  res, err := resource.New(context.Background(),
    resource.WithOS(),
    resource.WithTelemetrySDK(),
    resource.WithHost(),
    resource.WithAttributes(
      parameter.Resources...,
    ),
  )
  if err != nil {
    return nil, err
  }

  res, err = resource.Merge(resource.Default(), res)
  if err != nil {
    return nil, err
  }

  // Tracer
  traceProvider, err := newTraceProvider(context.Background(), res, parameter.Options...)
  if err != nil {
    return nil, err
  }
  shutdowns = append(shutdowns, traceProvider.Shutdown)
  otel.SetTracerProvider(traceProvider)

  metricProvider, err := newMetricProvider(res, parameter.Collectors...)
  if err != nil {
    return nil, err
  }
  shutdowns = append(shutdowns, metricProvider.Shutdown)
  otel.SetMeterProvider(metricProvider)

  // Propagator
  propagator := newPropagator()
  otel.SetTextMapPropagator(propagator)

  result := func(ctx context.Context) error {
    var err error
    for _, shutdown := range shutdowns {
      if shutdown != nil {
        err2 := shutdown(ctx)
        if err2 != nil {
          errors.Join(err, err2)
        }
      }
    }
    return err
  }

  return result, nil
}

func GetFromKafkaMessage(parent context.Context, message *sarama.ConsumerMessage) context.Context {
  return otel.GetTextMapPropagator().
    Extract(parent, otelsarama.NewConsumerMessageCarrier(message))
}

func Extract(metadata *types.EventMetadata, msg *sarama.ConsumerMessage) context.Context {
  return GetFromKafkaMessage(metadata.AsContext(context.Background()), msg)
}

func Inject(ctx context.Context, message *sarama.ProducerMessage) {
  otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(message))
}
