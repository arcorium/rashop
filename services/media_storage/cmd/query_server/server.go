package main

import (
  "context"
  "errors"
  "github.com/arcorium/rashop/services/media_storage/config"
  "github.com/arcorium/rashop/services/media_storage/constant"
  "github.com/arcorium/rashop/services/media_storage/internal/api/grpc/handler"
  "github.com/arcorium/rashop/services/media_storage/internal/app/query"
  "github.com/arcorium/rashop/services/media_storage/internal/app/service"
  "github.com/arcorium/rashop/services/media_storage/internal/infra/client"
  "github.com/arcorium/rashop/services/media_storage/internal/infra/model"
  "github.com/arcorium/rashop/services/media_storage/internal/infra/persistence"
  sharedConf "github.com/arcorium/rashop/shared/config"
  "github.com/arcorium/rashop/shared/database"
  "github.com/arcorium/rashop/shared/grpc/interceptor/log"
  "github.com/arcorium/rashop/shared/interfaces"
  "github.com/arcorium/rashop/shared/logger"
  otelUtil "github.com/arcorium/rashop/shared/otel"
  "github.com/arcorium/rashop/shared/types"
  sharedUtil "github.com/arcorium/rashop/shared/util"
  promProv "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "github.com/uptrace/bun"
  "github.com/uptrace/bun/extra/bunotel"
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

const serviceName = constant.ServiceQueryName

func NewServer(serverConfig *config.QueryServer) (*Server, error) {
  svr := &Server{
    ServerBase: interfaces.NewServer(),
    config:     serverConfig,
  }

  err := svr.setup()
  return svr, err
}

type Server struct {
  interfaces.ServerBase

  config *config.QueryServer
  db     *bun.DB

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
  mediaRepo := persistence.NewMediaPostgres(s.db)

  // Client
  storageClient, err := client.NewMinIOStorageClient(false, constant.StorageBucket, client.MinIOStorageConfig{
    Address:         s.config.Storage.Address,
    AccessKeyID:     s.config.Storage.AccessKey,
    SecretAccessKey: s.config.Storage.SecretKey,
  })
  if err != nil {
    return err
  }

  // Service
  queryConfig := service.DefaultMediaQueryConfig(query.CommonHandlerParameter{
    Persistent: mediaRepo,
    Storage:    storageClient,
  })
  querySvc := service.NewMediaQuery(queryConfig)

  // Handler
  // gRPC
  queryHandler := handler.NewMediaQuery(querySvc)
  queryHandler.Register(s.grpcServer)

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
