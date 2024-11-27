package producer

import (
	"context"
	"fmt"
	config "github.com/hyperfyodor/yq_sample/internal/config/producer"
	"github.com/hyperfyodor/yq_sample/internal/grpc"
	"github.com/hyperfyodor/yq_sample/internal/helpers"
	metricssrvr "github.com/hyperfyodor/yq_sample/internal/metrics"
	metrics "github.com/hyperfyodor/yq_sample/internal/metrics/producer"
	"github.com/hyperfyodor/yq_sample/internal/profiling"
	service "github.com/hyperfyodor/yq_sample/internal/service/producer"
	"github.com/hyperfyodor/yq_sample/internal/storage"
	"log"
	"log/slog"
	"strconv"
	"sync"
)

type App struct {
	log             *slog.Logger
	config          *config.Config
	producerService *service.Producer
	postgres        *storage.PostgresStorage
	numberOfWorkers int
}

func MustLoad(ctx context.Context) *App {
	cfg := config.MustLoad()
	producerMetrics := metrics.MustLoad()
	connection := helpers.ConnectionString(
		cfg.Db.Username,
		cfg.Db.Password,
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.Name,
		cfg.Db.SslMode,
	) + fmt.Sprintf("&pool_max_conns=%v", cfg.Db.PoolSize)

	poolSize, err := strconv.Atoi(cfg.Db.PoolSize)

	if err != nil {
		log.Printf("failed to get pool size: %v", err)
		panic(err)
	}

	numberOfWorkers := poolSize * 64

	postgres, err := storage.NewPostgresStorage(ctx, connection, true)

	if err != nil {
		log.Printf("failed to setup a storage: %v", err)
		panic(err)
	}

	publisher, err := grpc.NewGrpcPublisher(cfg.GrpcServer.Host, cfg.GrpcServer.Port, cfg.Mps)

	if err != nil {
		log.Printf("failed to setup a publisher: %v", err)
		panic(err)
	}

	logger := helpers.SetupLogger(cfg.LoggingLevel, cfg.LoggingType)

	producer := service.New(
		logger,
		postgres,
		publisher,
		postgres,
		producerMetrics,
	)

	return &App{
		log:             logger,
		config:          cfg,
		producerService: producer,
		numberOfWorkers: numberOfWorkers,
		postgres:        postgres,
	}
}

func (app *App) Start(ctx context.Context) {
	app.log.Info("Starting producer")

	jobs := make(chan int, app.numberOfWorkers)
	errors := make(chan error, app.numberOfWorkers)
	wg := sync.WaitGroup{}
	wg.Add(app.numberOfWorkers)

	for i := range app.numberOfWorkers {
		go app.worker(i, jobs, errors, ctx, &wg)
	}

L:
	for i := range app.config.MaxBacklog {
		jobs <- i
		select {
		case <-ctx.Done():
			close(jobs)
			break L
		case err := <-errors:
			app.log.Error("one of workers failed to produce task, shutting down", helpers.SlErr(err))
			close(jobs)
			break L
		default:

		}
	}

	app.log.Info("waiting for workers to finish")
	wg.Wait()
	app.log.Info("all workers finished")
	app.postgres.Close()
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

func (app *App) worker(id int, jobs <-chan int, errs chan<- error, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := app.log.With("worker_id", id)
	for range jobs {
		id, err := app.producerService.Produce(ctx)

		if err != nil {
			logger.Error("failed to produce task", helpers.SlErr(err))
			errs <- err
		}

		logger.Debug("successfully produced task", slog.String("id", strconv.Itoa(id)))
	}
}
