package main

import (
  "context"
  "errors"
  "github.com/IBM/sarama"
  sharedConf "github.com/arcorium/rashop/shared/config"
  "github.com/arcorium/rashop/shared/database"
  "github.com/arcorium/rashop/shared/grpc/interceptor/log"
  "github.com/arcorium/rashop/shared/logger"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  promProv "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/collectors"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "github.com/uptrace/bun"
  "github.com/uptrace/bun/extra/bunotel"
  "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
  "go.opentelemetry.io/otel/propagation"
  "go.opentelemetry.io/otel/sdk/resource"
  sdktrace "go.opentelemetry.io/otel/sdk/trace"
  semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/grpc/health"
  "google.golang.org/grpc/health/grpc_health_v1"
  "google.golang.org/grpc/reflection"
  "mini-shop/services/user/config"
  "mini-shop/services/user/constant"
  "mini-shop/services/user/internal/api/grpc/handler"
  "mini-shop/services/user/internal/api/messaging/consumer"
  "mini-shop/services/user/internal/api/messaging/dispatcher"
  "mini-shop/services/user/internal/app/service"
  "mini-shop/services/user/internal/infra/model"
  "mini-shop/services/user/internal/infra/persistence/pg"
  "net"
  "net/http"
  "os"
  "os/signal"
  "sync"
  "syscall"
)

func NewServer(serverConfig *config.QueryServer) (*Server, error) {
  svr := &Server{
    config: serverConfig,
  }

  err := svr.setup()
  return svr, err
}

type Server struct {
  config   *config.QueryServer
  db       *bun.DB
  consumer sarama.ConsumerGroup // It should be sarama.Consumer, because each instance should handle it all

  grpcServer     *grpc.Server
  metricServer   *http.Server
  logger         logger.ILogger
  exporter       sdktrace.SpanExporter
  tracerProvider *sdktrace.TracerProvider

  wg sync.WaitGroup
}

func (s *Server) validationSetup() {
  validator := sharedUtil.GetValidator()
  types.RegisterDefaultNullableValidations(validator)
}

func (s *Server) setupOtel() (*promProv.ServerMetrics, *prometheus.Registry, error) {
  var err error
  // Metrics
  metrics := promProv.NewServerMetrics(
    promProv.WithServerHandlingTimeHistogram(
      promProv.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
    ),
  )

  reg := prometheus.NewRegistry()
  reg.MustRegister(metrics)

  // Trace
  // Exporter
  s.exporter, err = otlptracegrpc.New(context.Background(),
    otlptracegrpc.WithInsecure(),
    otlptracegrpc.WithEndpoint(s.config.OTLPGRPCCollectorAddress),
  )
  if err != nil {
    return nil, nil, err
  }

  // Resource
  res, err := resource.New(context.Background(),
    resource.WithAttributes(
      semconv.ServiceName(constant.SERVICE_NAME),
      semconv.ServiceVersion(constant.SERVICE_VERSION),
    ))

  bsp := sdktrace.NewBatchSpanProcessor(s.exporter)
  s.tracerProvider = sdktrace.NewTracerProvider(
    sdktrace.WithSampler(sdktrace.AlwaysSample()),
    sdktrace.WithSpanProcessor(bsp),
    sdktrace.WithResource(res),
  )

  otel.SetTracerProvider(s.tracerProvider)
  otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
    propagation.TraceContext{}, propagation.Baggage{},
  ))

  return metrics, reg, nil
}

func (s *Server) grpcServerSetup() error {
  // Log
  s.logger = logger.GetGlobal()
  zaplogger, ok := s.logger.(*logger.ZapLogger)
  if !ok {
    return errors.New("logger is not of expected type, expected zap")
  }
  zapLogger := log.ZapLogger(zaplogger.Internal)

  metrics, reg, err := s.setupOtel()
  if err != nil {
    return err
  }

  exemplarFromCtx := func(ctx context.Context) prometheus.Labels {
    if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
      return prometheus.Labels{"traceID": span.TraceID().String()}
    }
    return nil
  }
  s.grpcServer = grpc.NewServer(
    grpc.StatsHandler(otelgrpc.NewServerHandler()), // tracing
    grpc.ChainUnaryInterceptor(
      recovery.UnaryServerInterceptor(),
      logging.UnaryServerInterceptor(zapLogger), // logging
    ),
    grpc.ChainStreamInterceptor(
      recovery.StreamServerInterceptor(),
      logging.StreamServerInterceptor(zapLogger), // logging
      metrics.StreamServerInterceptor(promProv.WithExemplarFromContext(exemplarFromCtx)),
    ),
  )

  if sharedConf.IsDebug() {
    reflection.Register(s.grpcServer)
    reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
    reg.MustRegister(collectors.NewGoCollector())
  }

  metrics.InitializeMetrics(s.grpcServer)
  // Metric endpoint
  mux := http.NewServeMux()
  mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
    Registry:          reg,
    EnableOpenMetrics: true,
  }))

  s.metricServer = &http.Server{Handler: mux, Addr: s.config.MetricAddress()}

  return nil
}

func (s *Server) databaseSetup() error {
  // Database
  var err error
  s.db, err = database.OpenPostgresWithConfig(&s.config.PostgresDatabase, sharedConf.IsDebug())
  if err != nil {
    return err
  }
  // Add trace hook
  s.db.AddQueryHook(bunotel.NewQueryHook(
    bunotel.WithFormattedQueries(true),
  ))
  model.RegisterBunModels(s.db)

  return nil
}

func (s *Server) createKafkaConfig() (*sarama.Config, error) {
  conf := sarama.NewConfig()
  if s.config.KafkaVersion != "" {
    var err error
    conf.Version, err = sarama.ParseKafkaVersion(s.config.KafkaVersion)
    if err != nil {
      return nil, err
    }
  } else {
    conf.Version = constant.DefaultKafkaVersion
  }
  conf.Producer.RequiredAcks = sarama.WaitForAll
  conf.Producer.Return.Successes = true
  return conf, nil
}

func (s *Server) consumerSetup() error {
  consumerCfg, err := s.createKafkaConfig()
  if err != nil {
    return err
  }

  group, err := sarama.NewConsumerGroup(s.config.Addresses, s.config.GroupId, consumerCfg)
  if err != nil {
    return err
  }

  s.consumer = group
  return err
}

func (s *Server) setup() error {
  s.validationSetup()

  if err := s.grpcServerSetup(); err != nil {
    return err
  }
  if err := s.databaseSetup(); err != nil {
    return err
  }

  // Repository
  // Persistent
  custRepo := pg.NewCustomer(s.db)

  // Publisher
  //kafkaPublisher := publisher.NewKafka(s.producer, serde.JsonSerializer{})

  // Service
  queryConsumerConfig := service.DefaultCustomerQueryConsumerConfig(custRepo)
  queryConsumerSvc := service.NewCustomerQueryConsumer(queryConsumerConfig)
  queryConfig := service.DefaultCustomerQueryConfig(custRepo)
  querySvc := service.NewCustomerQuery(queryConfig)

  // Handler
  customerQueryHandler := consumer.NewCustomerQueryHandler(queryConsumerSvc)
  // Consumer
  err := s.consumerSetup()
  if err != nil {
    return err
  }
  consumerDispatcher := dispatcher.NewQueryConsumerGroup(100, customerQueryHandler, serde.JsonAnyDeserializer{})
  dispatcher.RunQueryConsumerGroup(s.consumer, consumerDispatcher, constant.CustomerQuerySubscribedDomainEvent...)

  // gRPC
  queryHandler := handler.NewCustomerQuery(querySvc)
  queryHandler.Register(s.grpcServer)

  // Health check
  healthHandler := health.NewServer()
  grpc_health_v1.RegisterHealthServer(s.grpcServer, healthHandler)
  healthHandler.SetServingStatus(constant.SERVICE_NAME, grpc_health_v1.HealthCheckResponse_SERVING)

  return nil
}

func (s *Server) shutdown() {
  ctx := context.Background()

  s.grpcServer.GracefulStop()
  s.metricServer.Shutdown(ctx)
  s.wg.Wait()

  // OTEL
  err := s.exporter.Shutdown(ctx)
  if err != nil {
    logger.Warn(err.Error())
  }

  if err := s.tracerProvider.Shutdown(ctx); err != nil {
    logger.Warn(err.Error())
  }

  // Database
  if err := s.db.Close(); err != nil {
    logger.Warn(err.Error())
  }

  // Kafka

  if err := s.consumer.Close(); err != nil {
    logger.Warn(err.Error())
  }

  logger.Info("QueryServer Stopped!")
}

func (s *Server) Run() error {
  listener, err := net.Listen("tcp", s.config.Address())
  if err != nil {
    return err
  }

  // Run grpc server
  go func() {
    s.wg.Add(1)
    defer s.wg.Done()

    logger.Infof("%s listening on %s", constant.SERVICE_QUERY_NAME, s.config.Address())

    err = s.grpcServer.Serve(listener)
    logger.Infof("%s stopping", constant.SERVICE_QUERY_NAME)
    if err != nil && !errors.Is(err, http.ErrServerClosed) {
      logger.Warnf("%s failed to serve: %s", constant.SERVICE_QUERY_NAME, err)
    }
  }()

  go func() {
    s.wg.Add(1)
    defer s.wg.Done()
    logger.Infof("Metrics %s listening on %s", constant.SERVICE_QUERY_NAME, s.config.MetricAddress())

    err = s.metricServer.ListenAndServe()
    logger.Infof("Metrics %s stopping", constant.SERVICE_QUERY_NAME)
    if err != nil && !errors.Is(err, http.ErrServerClosed) {
      logger.Warnf("%s failed to serve: %s", constant.SERVICE_QUERY_NAME, err)
    }
  }()

  quitChan := make(chan os.Signal, 1)
  defer close(quitChan)

  signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
  <-quitChan

  s.shutdown()
  return err
}
