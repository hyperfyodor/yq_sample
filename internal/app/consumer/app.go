package consumer

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	config "github.com/hyperfyodor/yq_sample/internal/config/consumer"
	consumersrvr "github.com/hyperfyodor/yq_sample/internal/grpc"
	"github.com/hyperfyodor/yq_sample/internal/helpers"
	metricssrvr "github.com/hyperfyodor/yq_sample/internal/metrics"
	metrics "github.com/hyperfyodor/yq_sample/internal/metrics/consumer"
	"github.com/hyperfyodor/yq_sample/internal/profiling"
	service "github.com/hyperfyodor/yq_sample/internal/service/consumer"
	"github.com/hyperfyodor/yq_sample/internal/storage"
	"github.com/hyperfyodor/yq_sample/proto/consumer/gen"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	config     *config.Config
	grpcServer *grpc.Server
	postgres   *storage.PostgresStorage
}

func MustLoad(ctx context.Context) *App {
	cfg := config.MustLoad()
	consumerMetrics := metrics.MustLoad()
	connection := helpers.ConnectionString(
		cfg.Db.Username,
		cfg.Db.Password,
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.Name,
		cfg.Db.SslMode,
	) + fmt.Sprintf("&pool_max_conns=%v", cfg.Db.PoolSize)

	postgres, err := storage.NewPostgresStorage(ctx, connection, true)

	if err != nil {
		log.Printf("failed to setup a storage: %v", err)
		panic(err)
	}

	logger := helpers.SetupLogger(cfg.LoggingLevel, cfg.LoggingType)

	consumer := service.New(
		logger,
		postgres,
		consumerMetrics,
	)

	limiter := rate.NewLimiter(rate.Limit(cfg.Mcr), 1)

	server := consumersrvr.NewConsumerServer(consumer, limiter)

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(interceptorLogger(logger), loggingOpts...),
	))

	gen.RegisterConsumerServiceServer(grpcServer, server)

	return &App{
		log:        logger,
		config:     cfg,
		grpcServer: grpcServer,
		postgres:   postgres,
	}
}

func (app *App) Start() {
	lis, err := net.Listen("tcp", ":"+app.config.GrpcServer.Port)
	if err != nil {
		panic(err)
	}

	app.log.Info("grpc server listening at", "addr", lis.Addr())
	if err := app.grpcServer.Serve(lis); err != nil {
		app.log.Error("failed to serve", helpers.SlErr(err))

		panic(err)
	}
}

func (app *App) StartMetrics() {
	app.log.Info("starting metrics", slog.String("addr", ":"+app.config.MetricsPort+"/metrics"))
	if err := metricssrvr.Listen(app.config.MetricsPort); err != nil {
		app.log.Error("failed to start metrics server", slog.String("addr", ":"+app.config.MetricsPort+"/metrics"))
	}
}

func (app *App) StartProfiling() {
	app.log.Info("starting profiling", slog.String("addr", ":"+app.config.ProfilingPort))
	if err := profiling.Listen(app.config.ProfilingPort); err != nil {
		app.log.Error("failed to start profiling server", slog.String("addr", ":"+app.config.ProfilingPort))
	}
}

func (app *App) Stop() {
	app.grpcServer.GracefulStop()
	app.postgres.Close()
	app.log.Info("grpc server stopped")
}

func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
