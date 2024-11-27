package producer

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
)

var taskSaved bool
var taskPublished bool
var totalProduced int

type stub struct{}

func (*stub) SaveTask(ctx context.Context, taskType int, taskValue int) (int, error) {
	taskSaved = true
	return 1, nil
}

func (*stub) State(ctx context.Context, id int) (string, error) {
	return "", nil
}

func (*stub) Publish(ctx context.Context, taskId int, taskType int, taskValue int) error {
	taskPublished = true
	return nil
}

func (*stub) TotalProducedInc() {
	totalProduced++
}

func TestProducerHappyPath(t *testing.T) {
	ctx := context.Background()
	m := &stub{}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := New(log, m, m, m, m)

	id, err := service.Produce(ctx)

	assert.Equal(t, 1, id)
	assert.NoError(t, err)
	assert.True(t, taskSaved)
	assert.True(t, taskPublished)
	assert.Equal(t, 1, totalProduced)
}
