package main

import (
  "context"
  "errors"
  "github.com/IBM/sarama"
  "github.com/arcorium/rashop/services/token/config"
  "github.com/arcorium/rashop/services/token/constant"
  "github.com/arcorium/rashop/services/token/internal/api/grpc/handler"
  "github.com/arcorium/rashop/services/token/internal/api/messaging/consumer"
  "github.com/arcorium/rashop/services/token/internal/api/messaging/dispatcher"
  "github.com/arcorium/rashop/services/token/internal/app/command"
  commandCon "github.com/arcorium/rashop/services/token/internal/app/command/consumer"
  "github.com/arcorium/rashop/services/token/internal/app/service"
  "github.com/arcorium/rashop/services/token/internal/infra/persistence"
  "github.com/arcorium/rashop/services/token/internal/infra/publisher"
  sharedConf "github.com/arcorium/rashop/shared/config"
  "github.com/arcorium/rashop/shared/grpc/interceptor/log"
  "github.com/arcorium/rashop/shared/interfaces"
  "github.com/arcorium/rashop/shared/logger"
  "github.com/arcorium/rashop/shared/messaging/kafka"
  otelUtil "github.com/arcorium/rashop/shared/otel"
  "github.com/arcorium/rashop/shared/serde"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  promProv "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "github.com/redis/go-redis/extra/redisotel/v9"
  "github.com/redis/go-redis/v9"
  "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
  semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/grpc/health"
  "google.golang.org/grpc/health/grpc_health_v1"
  "google.golang.org/grpc/reflection"
  "net"
  "net/http"
  "os"
  "os/signal"
  "sync"
  "syscall"
)

const serviceName = constant.ServiceCommandName

func NewServer(serverConfig *config.CommandServer) (*Server, error) {
  svr := &Server{
    ServerBase: interfaces.NewServer(),
    config:     serverConfig,
  }

  err := svr.setup()
  return svr, err
}

type Server struct {
  interfaces.ServerBase

  config   *config.CommandServer
  db       *redis.Client
  producer sarama.SyncProducer
  consumer sarama.ConsumerGroup

  grpcServer   *grpc.Server
  metricServer *http.Server
  shutdownFunc otelUtil.ShutdownFunc

  wg sync.WaitGroup
}

func (s *Server) validationSetup() {
  validator := sharedUtil.GetValidator()
  types.RegisterDefaultNullableValidations(validator)
}

func (s *Server) grpcServerSetup() (*promProv.ServerMetrics, error) {
  // Log
  zaplogger, ok := logger.GetGlobal().(*logger.ZapLogger)
  if !ok {
    return nil, errors.New("logger is not of expected type, expected zap")
  }
  zapLogger := log.ZapLogger(zaplogger.Internal)

  // OTEL
  metrics := promProv.NewServerMetrics(
    promProv.WithServerHandlingTimeHistogram(
      promProv.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
    ),
  )

  shutdownFunc, err := otelUtil.Setup(
    otelUtil.SetupParameter{
      Resources: []attribute.KeyValue{
        semconv.ServiceName(constant.ServiceName),
        semconv.ServiceVersion(constant.ServiceVersion),
        semconv.ServiceInstanceID(s.Identity())},
      Options: []otlptracegrpc.Option{
        otlptracegrpc.WithInsecure(),
        otlptracegrpc.WithEndpoint(s.config.OTLPGRPCCollectorAddress)},
      Collectors: []prometheus.Collector{
        metrics,
      },
    },
  )

  if err != nil {
    return nil, err
  }
  s.shutdownFunc = shutdownFunc

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
      metrics.UnaryServerInterceptor(promProv.WithExemplarFromContext(exemplarFromCtx)),
    ),
    grpc.ChainStreamInterceptor(
      recovery.StreamServerInterceptor(),
      logging.StreamServerInterceptor(zapLogger), // logging
      metrics.StreamServerInterceptor(promProv.WithExemplarFromContext(exemplarFromCtx)),
    ),
  )

  return metrics, nil
}

func (s *Server) databaseSetup() error {
  // Database
  s.db = redis.NewClient(&redis.Options{
    Addr:     s.config.RedisAddress,
    Username: s.config.RedisUsername,
    Password: s.config.RedisPassword,
  })

  // Check connection
  if err := s.db.Ping(context.Background()).Err(); err != nil {
    return err
  }

  // Add tracing and metrics
  if err := redisotel.InstrumentTracing(s.db); err != nil {
    return err
  }

  if err := redisotel.InstrumentMetrics(s.db); err != nil {
    return err
  }

  return nil
}

func (s *Server) consumerSetup() error {
  conf := kafka.DefaultConfig(s.config.Broker.KafkaVersion,
    constant.DefaultKafkaVersion,
    kafka.WithDefaultConsumerGroup(s.Identity()))

  group, err := kafka.DefaultSyncGroupConsumer(s.config.Broker.Addresses, s.config.Broker.GroupId, conf)
  if err != nil {
    return err
  }

  s.consumer = group
  return err
}

func (s *Server) publisherSetup() error {
  cfg := kafka.DefaultConfig(s.config.Broker.KafkaVersion,
    constant.DefaultKafkaVersion,
    kafka.WithDefaultProducer(),
  )

  producer, err := kafka.DefaultSyncProducer(s.config.Broker.Addresses, cfg, kafka.WithOTELSyncProducer())
  if err != nil {
    return err
  }

  s.producer = producer
  return nil
}

func (s *Server) setup() error {
  s.validationSetup()

  metrics, err := s.grpcServerSetup()
  if err != nil {
    return err
  }
  if err := s.databaseSetup(); err != nil {
    return err
  }

  // Repository
  // Persistent
  tokenRepo := persistence.NewRedisToken(s.db, nil) // Use default config

  // Publisher
  err = s.publisherSetup()
  if err != nil {
    return err
  }

  kafkaPublisher := publisher.NewKafka(publisher.KafkaTopic{
    DomainEvent:      constant.TokenDomainEventTopic,
    IntegrationEvent: constant.TokenIntegrationEventTopic,
  }, s.producer, serde.JsonSerializer{})

  // DQL
  dlqPublisher := kafka.NewForwarder(constant.TokenDlqTopic, s.producer)

  // Service
  var cfg command.TokenGenerationConfig
  if s.config.Duration.IsSingle() {
    cfg = command.NewSingleTokenGenerationConfig(s.config.Duration.SingleExpirationTime)
  } else {
    cfg.VerificationTokenExpiration = s.config.Duration.VerificationTokenExpiryTime
    cfg.ResetTokenExpiration = s.config.Duration.ResetTokenExpiryTime
    cfg.GeneralTokenExpiration = s.config.Duration.GeneralTokenExpiryTime
    cfg.LoginTokenExpiration = s.config.Duration.LoginTokenExpiryTime
  }

  commandFactory := service.DefaultTokenCommandFactory(command.CommonHandlerParameter{
    Persistent: tokenRepo,
    Publisher:  kafkaPublisher,
  }, cfg)
  commandSvc := service.NewTokenCommand(commandFactory)

  commandConsumerFactory := service.DefaultTokenCommandConsumerFactory(commandCon.CommonHandlerParameter{
    Persistent: tokenRepo,
    Generate:   commandFactory.Generate,
  })
  commandConsumerSvc := service.NewTokenCommandConsumer(commandConsumerFactory)

  // Handler
  // gRPC
  commandHandler := handler.NewTokenCommand(commandSvc)
  commandHandler.Register(s.grpcServer)

  // Messaging
  commandConsumerHandler := consumer.NewTokenCommand(commandConsumerSvc)
  if err := s.consumerSetup(); err != nil {
    return err
  }

  consumerDispatcher := dispatcher.NewTokenCommandConsumerGroup(dlqPublisher, commandConsumerHandler, serde.JsonAnyDeserializer{})
  consumerDispatcher.Run(context.Background(), s.consumer, constant.ListeningTopics...)

  // Health check
  healthHandler := health.NewServer()
  grpc_health_v1.RegisterHealthServer(s.grpcServer, healthHandler)
  healthHandler.SetServingStatus(constant.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
  healthHandler.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_SERVING)

  // Reflection
  if sharedConf.IsDebug() {
    reflection.Register(s.grpcServer)
  }

  // Metric HTTP
  metrics.InitializeMetrics(s.grpcServer)
  mux := http.NewServeMux()
  mux.Handle("/metrics", promhttp.Handler())
  s.metricServer = &http.Server{Handler: mux, Addr: s.config.MetricAddress()}

  return nil
}

func (s *Server) shutdown() {
  ctx := context.Background()

  s.grpcServer.GracefulStop()
  s.metricServer.Shutdown(ctx)
  s.wg.Wait()

  // OTEL
  if err := s.shutdownFunc(ctx); err != nil {
    logger.Warn(err.Error())
  }

  // Database
  if err := s.db.Close(); err != nil {
    logger.Warn(err.Error())
  }

  // Messaging
  if err := s.producer.Close(); err != nil {
    logger.Warn(err.Error())
  }
  if err := s.consumer.Close(); err != nil {
    logger.Warn(err.Error())
  }

  logger.Infof("%s Stopped!", serviceName)
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

    logger.Infof("%s listening on %s", serviceName, s.config.Address())

    err = s.grpcServer.Serve(listener)
    logger.Infof("%s stopping", serviceName)
    if err != nil && !errors.Is(err, http.ErrServerClosed) {
      logger.Warnf("%s failed to serve: %s", serviceName, err)
    }
  }()

  go func() {
    s.wg.Add(1)
    defer s.wg.Done()
    logger.Infof("Metrics %s listening on %s", serviceName, s.config.MetricAddress())

    err = s.metricServer.ListenAndServe()
    logger.Infof("Metrics %s stopping", serviceName)
    if err != nil && !errors.Is(err, http.ErrServerClosed) {
      logger.Warnf("%s failed to serve: %s", serviceName, err)
    }
  }()

  quitChan := make(chan os.Signal, 1)
  defer close(quitChan)

  signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
  s.MarkRunning()
  <-quitChan

  s.shutdown()
  return err
}
