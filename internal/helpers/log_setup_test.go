package helpers

import (
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestLoggingLevel(t *testing.T) {
	assert.Equal(t, loggingLevel("debug"), slog.LevelDebug)
	assert.Equal(t, loggingLevel("info"), slog.LevelInfo)
	assert.Equal(t, loggingLevel("warn"), slog.LevelWarn)
	assert.Equal(t, loggingLevel("error"), slog.LevelError)
	assert.Equal(t, loggingLevel("none"), slog.LevelDebug)
}

func TestSetupTextLogger(t *testing.T) {
	testLogger := SetupLogger("none", "text")
	assert.NotNil(t, testLogger)
	assert.IsTypef(t, &slog.TextHandler{}, testLogger.Handler(), "logger should be of type slog.TextHandler")
}

func TestSetupJSONLogger(t *testing.T) {
	testLogger := SetupLogger("none", "json")
	assert.NotNil(t, testLogger)
	assert.IsTypef(t, &slog.JSONHandler{}, testLogger.Handler(), "logger should be of type slog.JSONHandler")
}

func TestSetupDefaultLogger(t *testing.T) {
	testLogger := SetupLogger("none", "none")
	assert.NotNil(t, testLogger)
	assert.IsTypef(t, &slog.TextHandler{}, testLogger.Handler(), "logger should be of type slog.TextHandler")
}
