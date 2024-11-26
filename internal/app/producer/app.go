package producer

import (
	"context"
	"fmt"
	config "github.com/hyperfyodor/yq_sample/internal/config/producer"
	"github.com/hyperfyodor/yq_sample/internal/grpc"
	"github.com/hyperfyodor/yq_sample/internal/helpers"
	metrics "github.com/hyperfyodor/yq_sample/internal/metrics/producer"
	service "github.com/hyperfyodor/yq_sample/internal/service/producer"
	"github.com/hyperfyodor/yq_sample/internal/storage"
	"log"
	"log/slog"
	"strconv"
	"sync"
	"time"
)

type App struct {
	log             *slog.Logger
	config          *config.Config
	producerService *service.Producer
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

	numberOfWorkers := poolSize * 4

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
	}
}

func (app *App) Start(ctx context.Context) {
	tick := time.Tick(10 * time.Millisecond)

	jobs := make(chan int, app.numberOfWorkers)
	errors := make(chan error, app.numberOfWorkers)
	wg := sync.WaitGroup{}
	wg.Add(app.numberOfWorkers)

	for i := range app.numberOfWorkers {
		go app.worker(i, jobs, errors, ctx, &wg)
	}

L:
	for i := range app.config.MaxBacklog {
		select {
		case <-ctx.Done():
			close(jobs)
			break L
		case err := <-errors:
			app.log.Error("one of workers failed to produce task, shutting down", helpers.SlErr(err))
			close(jobs)
			break L
		case <-tick:
			jobs <- i
		}
	}

	app.log.Info("waiting for workers to finish")
	wg.Wait()
	app.log.Info("all workers finished")
}

func (app *App) worker(id int, jobs <-chan int, errs chan<- error, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log := app.log.With("worker_id", id)
	for range jobs {
		id, err := app.producerService.Produce(ctx)

		if err != nil {
			log.Error("failed to produce task", helpers.SlErr(err))
			errs <- err
		}

		log.Debug("successfully produced task", slog.String("id", strconv.Itoa(id)))
	}
}
