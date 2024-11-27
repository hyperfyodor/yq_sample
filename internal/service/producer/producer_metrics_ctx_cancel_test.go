package producer

import (
	"context"
	metrics "github.com/hyperfyodor/yq_sample/internal/metrics/producer"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
)

func TestProducerMetricsCtxCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert := assert.New(t)
	s := &stub{}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	metrics := metrics.MustLoad()
	producer := New(log, s, s, s, metrics)

	id, err := producer.Produce(ctx)

	assert.Error(err)
	assert.Equal(0, id)
	assert.Equal(float64(0), testutil.ToFloat64(metrics.TotalProduced))

	metrics.Unregister()
}
