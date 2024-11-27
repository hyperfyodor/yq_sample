package consumer

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
)

func TestConsumerContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &stub{}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := New(log, s, s)

	cancel()
	err := service.Consume(ctx, 0, 0, 0)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
}
