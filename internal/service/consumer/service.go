package consumer

import (
	"context"
	"github.com/hyperfyodor/yq_sample/internal/helpers"
	"log/slog"
	"strconv"
	"time"
)

// StateUpdater changes task state somewhere
type StateUpdater interface {
	Done(ctx context.Context, id int) error
	Processing(ctx context.Context, id int) error
}

// MetricsForConsumer provides methods, that we call inside service, when the state of the task changes
type MetricsForConsumer interface {
	TaskJustReceived()
	TaskIsProcessing()
	TaskIsDone(taskType int, taskValue int)
}

type Consumer struct {
	log          *slog.Logger
	stateUpdater StateUpdater
	metrics      MetricsForConsumer
	counter      *typeCounter
}

func New(log *slog.Logger, stateUpdater StateUpdater, metrics MetricsForConsumer) *Consumer {
	if log == nil {
		panic("log is nil")
	}

	counter := newTypeCounter()

	return &Consumer{
		log:          log,
		stateUpdater: stateUpdater,
		metrics:      metrics,
		counter:      counter,
	}
}

func (consumer *Consumer) Consume(ctx context.Context, taskId, taskType, taskValue int) error {
	const op = "internal.service.consumer.Consume"

	log := consumer.log.With(
		slog.String("task_id", strconv.Itoa(taskId)),
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return helpers.WrapErr(op, ctx.Err())
	default:
		consumer.metrics.TaskJustReceived()

		log.Debug("setting task state to 'processing'")

		if err := consumer.stateUpdater.Processing(ctx, taskId); err != nil {
			log.Error("failed to set task state to 'processing'", helpers.SlErr(err))
			return helpers.WrapErr(op, err)
		}

		consumer.metrics.TaskIsProcessing()

		select {
		case <-time.After(time.Millisecond * time.Duration(taskValue)):

			log.Debug("setting task state to 'done'")

			if err := consumer.stateUpdater.Done(ctx, taskId); err != nil {
				log.Error("failed to set task state to 'done'", helpers.SlErr(err))
				return helpers.WrapErr(op, err)
			} else {
				consumer.metrics.TaskIsDone(taskType, taskValue)
				consumer.counter.Inc(taskType)

				log.Info("successfully consumed task",
					slog.String("id", strconv.Itoa(taskId)),
					slog.String("type", strconv.Itoa(taskType)),
					slog.String("value", strconv.Itoa(taskValue)),
					slog.String("total_type_count", strconv.Itoa(consumer.counter.Get(taskType))),
				)

			}
			return nil

		case <-ctx.Done():
			log.Debug("context canceled")
			return nil
		}

	}
}
