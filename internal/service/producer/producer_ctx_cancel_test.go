package producer

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
)

func TestProducerContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &stub{}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := New(log, s, s, s, s)

	cancel()
	id, err := service.Produce(ctx)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.True(t, errors.Is(err, context.Canceled))
}
