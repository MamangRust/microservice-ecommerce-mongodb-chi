package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/dotenv"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-pkg/middleware"
	otel_pkg "github.com/MamangRust/microservice-ecommerce-pkg/otel"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"github.com/grafana/pyroscope-go"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	Logger           logger.LoggerInterface
	MongoDB          *mongo.Database
	MongoClient      *mongo.Client
	Ctx              context.Context
	Cancel           context.CancelFunc
	CacheStore       *cache.CacheStore
	Redis            *redis.Client
	Telemetry        *otel_pkg.Telemetry
	Config           *Config
	RegisterServices func(*grpc.Server)
	PoolMetrics      *PoolMetrics

	// extraMu guards extraClients. Setiap bounded context punya server MongoDB
	// sendiri, jadi service yang perlu database milik context lain membuka
	// koneksi kedua — bukan memakai ulang MongoClient.
	extraMu      sync.Mutex
	extraClients map[string]*mongo.Client
}

func New(cfg *Config) (*GRPCServer, error) {
	if err := dotenv.Viper(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	if err := initPyroscope(cfg); err != nil {
		log.Printf("Warning: Failed to initialize pyroscope: %v", err)
	}

	telemetry := initTelemetry(cfg)
	if err := telemetry.Init(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize telemetry: %w", err)
	}

	cacheMetrics, err := observability.NewCacheMetrics("cache")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache metrics: %w", err)
	}

	l, err := logger.NewLogger(cfg.ServiceName, telemetry.GetLogger())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// DBName wajib diisi — fail-hard, bukan fallback ke database default.
	if cfg.DBName == "" {
		return nil, fmt.Errorf("DBName wajib diisi untuk service %q (lihat pkg/database/names.go)", cfg.ServiceName)
	}

	// Connect to MongoDB: koneksi (client) dipisah dari pemilihan database.
	// Setiap bounded context punya server sendiri ("1 database = 1 MongoDB"),
	// jadi service hanya membuka koneksi ke instance miliknya. Service yang
	// butuh context lain memakai ContextDatabase.
	mongoClient, err := database.NewMongoClientForContext(l, cfg.DBName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB for %s: %w", cfg.DBName, err)
	}
	mongoDB := database.NewMongoDatabase(mongoClient, cfg.DBName)

	ctx, cancel := context.WithCancel(context.Background())

	redisClient, err := initRedisServer(ctx, l, cfg.ServiceName)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	cacheStore := cache.NewCacheStore(redisClient, l, cacheMetrics)

	poolMetrics, err := NewPoolMetrics(cfg.ServiceName)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize pool metrics: %w", err)
	}

	return &GRPCServer{
		Logger:      l,
		MongoDB:     mongoDB,
		MongoClient: mongoClient,
		Ctx:         ctx,
		Cancel:      cancel,
		CacheStore:  cacheStore,
		Redis:       redisClient,
		Telemetry:   telemetry,
		Config:      cfg,
		PoolMetrics: poolMetrics,
	}, nil
}

func (s *GRPCServer) Run() error {
	defer s.Cleanup()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Config.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.Config.Port, err)
	}

	loadMonitor := resilience.NewLoadMonitor()
	circuitBreaker := resilience.NewCircuitBreaker(100, 10, s.Logger)
	requestLimiter := resilience.NewRequestLimiter(800, s.Logger)
	resilienceHandler := middleware.NewResilienceInterceptor(loadMonitor, circuitBreaker, requestLimiter)

	metricsInterceptor, err := middleware.NewMetricsInterceptor(s.Config.ServiceName)
	if err != nil {
		return fmt.Errorf("failed to initialize metrics interceptor: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(DefaultMaxConcurrentConn),
		grpc.InitialConnWindowSize(DefaultWindowSize),
		grpc.InitialWindowSize(DefaultWindowSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    DefaultKeepaliveTime,
			Timeout: DefaultKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             DefaultMinKeepaliveTime,
			PermitWithoutStream: true,
		}),
		grpc.ChainUnaryInterceptor(
			middleware.ContextMiddleware(30*time.Second, s.Logger),
			middleware.TraceUnaryServerInterceptor(s.Config.ServiceName),
			metricsInterceptor.UnaryInterceptor(),
			middleware.RecoveryMiddleware(s.Logger),
			middleware.PyroscopeUnaryInterceptor(),
			resilienceHandler.UnaryInterceptor(),
		),
	)

	if s.RegisterServices != nil {
		s.RegisterServices(grpcServer)
	}

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	if os.Getenv("ENABLE_REFLECTION") == "true" {
		reflection.Register(grpcServer)
		s.Logger.Info("gRPC reflection enabled")
	}

	monitoringDone := s.spawnMonitoringTask()
	cleanupDone := s.spawnCleanupTask()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	errChan := make(chan error, 1)
	go func() {
		s.Logger.Info("gRPC server starting",
			zap.Int("port", s.Config.Port),
			zap.String("address", lis.Addr().String()),
		)
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	select {
	case sig := <-sigChan:
		s.Logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	case err := <-errChan:
		s.Logger.Error("Server error", zap.Error(err))
		return err
	}

	return s.gracefulShutdown(grpcServer, healthServer, monitoringDone, cleanupDone)
}

func (s *GRPCServer) gracefulShutdown(
	grpcServer *grpc.Server,
	healthServer *health.Server,
	monitoringDone, cleanupDone <-chan struct{},
) error {
	s.Logger.Info("Starting graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer shutdownCancel()

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	s.Cancel()

	tasksDone := make(chan struct{})
	go func() {
		<-monitoringDone
		<-cleanupDone
		close(tasksDone)
	}()

	select {
	case <-tasksDone:
		s.Logger.Info("Background tasks stopped successfully")
	case <-shutdownCtx.Done():
		s.Logger.Warn("Background tasks shutdown timeout, forcing stop")
	}

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		s.Logger.Info("gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		s.Logger.Warn("Graceful shutdown timeout, forcing stop")
		grpcServer.Stop()
	}

	s.Logger.Info("Graceful shutdown completed")
	return nil
}

// ContextDatabase returns a handle to a database owned by another bounded
// context.
//
// Dipakai oleh service yang memang perlu koleksi lintas context (mis. order
// membaca product_stock_adjustments milik ec_catalog). Karena setiap context
// punya server MongoDB sendiri, ini membuka koneksi kedua lewat
// NewMongoClientForContext — bukan memakai ulang MongoClient milik service.
// Koneksi di-cache per context supaya tidak dibuka ulang tiap panggilan, dan
// ditutup semua di Cleanup.
func (s *GRPCServer) ContextDatabase(name string) (*mongo.Database, error) {
	if name == s.Config.DBName {
		return s.MongoDB, nil
	}

	s.extraMu.Lock()
	defer s.extraMu.Unlock()

	if s.extraClients == nil {
		s.extraClients = make(map[string]*mongo.Client, 1)
	}

	if client, ok := s.extraClients[name]; ok {
		return client.Database(name), nil
	}

	client, err := database.NewMongoClientForContext(s.Logger, name)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB for %s: %w", name, err)
	}

	s.extraClients[name] = client
	s.Logger.Info("Opened MongoDB connection for foreign context", zap.String("database", name))

	return client.Database(name), nil
}

func (s *GRPCServer) Cleanup() {
	s.Logger.Info("Cleaning up resources...")

	if s.MongoClient != nil {
		if err := s.MongoClient.Disconnect(context.Background()); err != nil {
			s.Logger.Error("Failed to close MongoDB connection", zap.Error(err))
		} else {
			s.Logger.Info("MongoDB connection closed")
		}
	}

	s.extraMu.Lock()
	for name, client := range s.extraClients {
		if err := client.Disconnect(context.Background()); err != nil {
			s.Logger.Error("Failed to close MongoDB connection",
				zap.String("database", name),
				zap.Error(err))
		} else {
			s.Logger.Info("MongoDB connection closed", zap.String("database", name))
		}
	}
	s.extraMu.Unlock()

	if s.Redis != nil {
		if err := s.Redis.Close(); err != nil {
			s.Logger.Error("Failed to close Redis connection", zap.Error(err))
		} else {
			s.Logger.Info("Redis connection closed")
		}
	}

	if s.Telemetry != nil {
		if err := s.Telemetry.Shutdown(context.Background()); err != nil {
			s.Logger.Error("Failed to shutdown telemetry", zap.Error(err))
		} else {
			s.Logger.Info("Telemetry shutdown successfully")
		}
	}

	s.Logger.Info("Cleanup completed")
}

func initPyroscope(cfg *Config) error {
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: cfg.ServiceName,
		ServerAddress:   os.Getenv("PYROSCOPE_SERVER"),
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
		Tags: map[string]string{
			"service": cfg.ServiceName,
			"env":     cfg.Environment,
			"version": cfg.ServiceVersion,
		},
	})
	return err
}

func initTelemetry(cfg *Config) *otel_pkg.Telemetry {
	samplingFraction := cfg.OtelSamplingFraction
	if samplingFraction <= 0 {
		samplingFraction = 1.0
	}

	return otel_pkg.NewTelemetry(otel_pkg.Config{
		ServiceName:            cfg.ServiceName,
		ServiceVersion:         cfg.ServiceVersion,
		Environment:            cfg.Environment,
		Endpoint:               cfg.OtelEndpoint,
		Insecure:               true,
		EnableRuntimeMetrics:   true,
		RuntimeMetricsInterval: 15 * time.Second,
		SamplingFraction:       samplingFraction,
	})
}

func initRedisServer(ctx context.Context, logger logger.LoggerInterface, serviceName string) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST"), viper.GetString("REDIS_PORT")),
		Password:     viper.GetString("REDIS_PASSWORD"),
		DB:           viper.GetInt("REDIS_DB"),
		DialTimeout:  RedisDialTimeout,
		ReadTimeout:  RedisReadTimeout,
		WriteTimeout: RedisWriteTimeout,
		PoolSize:     RedisPoolSize,
		MinIdleConns: RedisMinIdleConns,
	}), nil
}

func (s *GRPCServer) spawnMonitoringTask() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(MonitoringInterval)
		defer ticker.Stop()
		for {
			select {
			case <-s.Ctx.Done():
				return
			case <-ticker.C:
				s.monitorCache()
				s.monitorPoolMetrics()
			}
		}
	}()
	return done
}

func (s *GRPCServer) monitorPoolMetrics() {
	if s.PoolMetrics == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), MonitoringInterval)
	defer cancel()

	if s.Redis != nil {
		poolStats := s.Redis.PoolStats()
		s.PoolMetrics.RecordRedisStats(ctx,
			int64(poolStats.TotalConns),
			int64(poolStats.IdleConns),
			int64(poolStats.Hits),
			int64(poolStats.Misses),
		)
	}
}

func (s *GRPCServer) monitorCache() {
	refCount := s.CacheStore.GetRefCount()
	stats, err := s.CacheStore.GetStats(s.Ctx)
	if err != nil {
		s.Logger.Error("Failed to get cache stats", zap.Error(err))
		return
	}
	logLevel := zap.InfoLevel
	if refCount > CacheRefCountThreshold {
		logLevel = zap.WarnLevel
	}
	if ce := s.Logger.Check(logLevel, "Cache statistics"); ce != nil {
		ce.Write(
			zap.Int64("ref_count", refCount),
			zap.Int64("total_keys", stats.TotalKeys),
			zap.Float64("hit_rate", stats.HitRate),
			zap.String("memory_used", stats.MemoryUsedHuman),
		)
	}
}

func (s *GRPCServer) spawnCleanupTask() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(CleanupInterval)
		defer ticker.Stop()
		for {
			select {
			case <-s.Ctx.Done():
				return
			case <-ticker.C:
				s.CacheStore.ClearExpired(s.Ctx)
			}
		}
	}()
	return done
}
