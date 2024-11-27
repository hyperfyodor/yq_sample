package producer

import (
	"context"
	"github.com/hyperfyodor/yq_sample/internal/helpers"
	"log/slog"
	"math/rand/v2"
	"strconv"
)

// TaskSaver creates a task record with provided task type and task value, returns id
type TaskSaver interface {
	SaveTask(ctx context.Context, taskType int, taskValue int) (int, error)
}

// StateProvider returns task state for a provided task id
type StateProvider interface {
	State(ctx context.Context, id int) (string, error)
}

// TaskPublisher publishes a message with taskId, taskType and taskValue
type TaskPublisher interface {
	Publish(ctx context.Context, taskId int, taskType int, taskValue int) error
}

// MetricsForProducer provides a method to increase total produced counter
type MetricsForProducer interface {
	TotalProducedInc()
}

type Producer struct {
	log           *slog.Logger
	taskSaver     TaskSaver
	taskPublisher TaskPublisher
	stateProvider StateProvider
	metrics       MetricsForProducer
}

func New(
	log *slog.Logger,
	taskSaver TaskSaver,
	taskPublisher TaskPublisher,
	stateProvider StateProvider,
	metrics MetricsForProducer,
) *Producer {

	if log == nil {
		panic("log is nil")
	}

	return &Producer{
		log:           log,
		taskSaver:     taskSaver,
		taskPublisher: taskPublisher,
		stateProvider: stateProvider,
		metrics:       metrics,
	}
}

func (producer *Producer) Produce(ctx context.Context) (int, error) {
	const op = "internal.service.producer.Produce"

	log := producer.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		log.Debug("context canceled")
		return 0, helpers.WrapErr(op, ctx.Err())
	default:
		taskType := randRange(0, 10)
		taskValue := randRange(0, 100)

		id, err := producer.taskSaver.SaveTask(ctx, taskType, taskValue)

		if err != nil {
			log.Error("failed to save task", helpers.SlErr(err))
			return 0, helpers.WrapErr(op, err)
		}

		log.Debug("saved task", slog.String("task_id", strconv.Itoa(id)))

		err = producer.taskPublisher.Publish(ctx, id, taskType, taskValue)

		if err != nil {
			log.Error("failed to publish task", helpers.SlErr(err), slog.String("task_id", strconv.Itoa(id)))
			return 0, helpers.WrapErr(op, err)
		}

		producer.metrics.TotalProducedInc()
		log.Debug("task published", slog.String("task_id", strconv.Itoa(id)))
		return id, nil
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
