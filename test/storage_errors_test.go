package test

import (
	"github.com/hyperfyodor/yq_sample/test/suite/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageErrors(t *testing.T) {
	ctx, storageSuite := storage.New(t)

	id, err := storageSuite.Storage.SaveTask(ctx, 12, 0)

	assert.Error(t, err)
	assert.Equal(t, 0, id)

	id, err = storageSuite.Storage.SaveTask(ctx, 0, 1500)

	assert.Error(t, err)
	assert.Equal(t, 0, id)

	id, err = storageSuite.Storage.SaveTask(ctx, 1500, 1500)

	assert.Error(t, err)
	assert.Equal(t, 0, id)

	err = storageSuite.Storage.Done(ctx, 56)

	assert.Error(t, err)
	assert.Equal(t, 0, id)

	err = storageSuite.Storage.Processing(ctx, 5)

	assert.Error(t, err)
	assert.Equal(t, 0, id)

	state, err := storageSuite.Storage.State(ctx, 34)
	assert.Error(t, err)
	assert.Equal(t, "", state)

}
