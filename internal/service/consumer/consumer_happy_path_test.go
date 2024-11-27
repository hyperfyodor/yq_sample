package consumer

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
)

var processing bool
var done bool
var metricsReceived bool
var metricsProcessing bool
var metricsDone bool

type stub struct{}

func (*stub) TaskJustReceived()                      { metricsReceived = true }
func (*stub) TaskIsProcessing()                      { metricsProcessing = true }
func (*stub) TaskIsDone(taskType int, taskValue int) { metricsDone = true }

func (*stub) Done(ctx context.Context, id int) error       { done = true; return nil }
func (*stub) Processing(ctx context.Context, id int) error { processing = true; return nil }

func TestConsumerHappyPath(t *testing.T) {
	s := &stub{}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	consumer := New(log, s, s)

	err := consumer.Consume(context.Background(), 0, 0, 0)
	assert.NoError(t, err)
	assert.True(t, processing)
	assert.True(t, metricsReceived)
	assert.True(t, metricsProcessing)
	assert.True(t, metricsDone)
	assert.True(t, done)
}
