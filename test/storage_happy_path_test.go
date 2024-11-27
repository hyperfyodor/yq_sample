package test

import (
	"github.com/hyperfyodor/yq_sample/test/suite/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHappyPath(t *testing.T) {
	ctx, storageSuite := storage.New(t)

	id, err := storageSuite.Storage.SaveTask(ctx, 0, 0)
	assert.Nil(t, err)

	state, err := storageSuite.Storage.State(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "received", state)

	err = storageSuite.Storage.Processing(ctx, id)
	assert.Nil(t, err)

	state, err = storageSuite.Storage.State(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "processing", state)

	err = storageSuite.Storage.Done(ctx, id)
	assert.Nil(t, err)

	state, err = storageSuite.Storage.State(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "done", state)
}
