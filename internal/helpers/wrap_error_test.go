package helpers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorWrapping(t *testing.T) {
	testError := errors.New("test error")
	wrappedError := WrapErr("wrapped", testError)

	assert.Equal(t, errors.Is(wrappedError, testError), true)
}

func TestSlogErrorWrapping(t *testing.T) {
	testError := errors.New("test error")
	slErrorAttr := SlErr(testError)

	assert.Equal(t, slErrorAttr.Value.String(), "test error")
}
