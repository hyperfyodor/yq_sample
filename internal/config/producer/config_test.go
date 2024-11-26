package producer

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Tests if config overrides default values with values from env correctly
func TestLoadProducerConfigFromEnv(t *testing.T) {
	envVars := map[string]string{
		"PRD_DB_HOST":        "test_prd_host",
		"PRD_DB_PORT":        "test_prd_port",
		"PRD_DB_USERNAME":    "test_prd_username",
		"PRD_DB_PASSWORD":    "test_prd_password",
		"PRD_DB_NAME":        "test_prd_name",
		"PRD_DB_SSL_MODE":    "test_prd_ssl_mode",
		"PRD_DB_POOL_SIZE":   "test_prd_pool_size",
		"PRD_LOGGING_LEVEL":  "test_prd_logging_level",
		"PRD_LOGGING_TYPE":   "test_prd_logging_type",
		"PRD_GRPC_PORT":      "test_prd_grpc_port",
		"PRD_GRPC_HOST":      "test_prd_grpc_host",
		"PRD_METRICS_PORT":   "test_prd_metrics_port",
		"PRD_PROFILING_PORT": "test_prd_profiling_port",
		"PRD_MAX_BACKLOG":    "1",
		"PRD_MPS":            "1",
	}

	for k, v := range envVars {
		err := os.Setenv(k, v)
		if err != nil {
			t.Errorf("Failed to set env var: %s", k)
		}
	}

	config := MustLoad()

	assert.Equal(t, envVars["PRD_DB_HOST"], config.Db.Host)
	assert.Equal(t, envVars["PRD_DB_PORT"], config.Db.Port)
	assert.Equal(t, envVars["PRD_DB_USERNAME"], config.Db.Username)
	assert.Equal(t, envVars["PRD_DB_PASSWORD"], config.Db.Password)
	assert.Equal(t, envVars["PRD_DB_NAME"], config.Db.Name)
	assert.Equal(t, envVars["PRD_DB_SSL_MODE"], config.Db.SslMode)
	assert.Equal(t, envVars["PRD_DB_POOL_SIZE"], config.Db.PoolSize)
	assert.Equal(t, envVars["PRD_LOGGING_LEVEL"], config.LoggingLevel)
	assert.Equal(t, envVars["PRD_LOGGING_TYPE"], config.LoggingType)
	assert.Equal(t, envVars["PRD_GRPC_PORT"], config.GrpcServer.Port)
	assert.Equal(t, envVars["PRD_GRPC_HOST"], config.GrpcServer.Host)
	assert.Equal(t, 1, config.Mps)
	assert.Equal(t, 1, config.MaxBacklog)
	assert.Equal(t, envVars["PRD_METRICS_PORT"], config.MetricsPort)
	assert.Equal(t, envVars["PRD_PROFILING_PORT"], config.ProfilingPort)
}
