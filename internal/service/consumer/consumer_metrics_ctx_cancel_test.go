package consumer

import (
	"context"
	metrics "github.com/hyperfyodor/yq_sample/internal/metrics/consumer"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
)

func TestConsumerMetricsCtxCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert := assert.New(t)
	s := &stub{}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	metrics := metrics.MustLoad()
	consumer := New(log, s, metrics)

	err := consumer.Consume(ctx, 0, 0, 0)

	assert.Error(err)
	assert.Equal(float64(0), testutil.ToFloat64(metrics.TaskCountPerState.WithLabelValues("received")))
	assert.Equal(float64(0), testutil.ToFloat64(metrics.TaskCountPerState.WithLabelValues("processing")))
	assert.Equal(float64(0), testutil.ToFloat64(metrics.TaskCountPerState.WithLabelValues("done")))

	assert.Equal(float64(0), testutil.ToFloat64(metrics.TaskTotalReceived))
	assert.Equal(float64(0), testutil.ToFloat64(metrics.TaskTotalPerTaskType.WithLabelValues("0")))
	assert.Equal(float64(0), testutil.ToFloat64(metrics.TaskTotalValuePerTaskType.WithLabelValues("0")))

	metrics.Unregister()
}
