package consumer

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Tests if config overrides default values with values from env correctly
func TestLoadConsumerConfigFromEnv(t *testing.T) {
	envVars := map[string]string{
		"CSM_DB_HOST":        "test_csm_host",
		"CSM_DB_PORT":        "test_csm_port",
		"CSM_DB_USERNAME":    "test_csm_username",
		"CSM_DB_PASSWORD":    "test_csm_password",
		"CSM_DB_NAME":        "test_csm_name",
		"CSM_DB_POOL_SIZE":   "1",
		"CSM_DB_SSL_MODE":    "test_csm_ssl_mode",
		"CSM_LOGGING_LEVEL":  "test_csm_logging_level",
		"CSM_LOGGING_TYPE":   "test_csm_logging_type",
		"CSM_GRPC_PORT":      "test_csm_grpc_port",
		"CSM_MCR":            "1",
		"CSM_METRICS_PORT":   "test_csm_metrics_port",
		"CSM_PROFILING_PORT": "test_csm_profiling_port",
	}

	for k, v := range envVars {
		err := os.Setenv(k, v)
		if err != nil {
			t.Errorf("Failed to set env var: %s", k)
		}
	}

	config := MustLoad()

	assert.Equal(t, envVars["CSM_DB_HOST"], config.Db.Host)
	assert.Equal(t, envVars["CSM_DB_PORT"], config.Db.Port)
	assert.Equal(t, envVars["CSM_DB_USERNAME"], config.Db.Username)
	assert.Equal(t, envVars["CSM_DB_PASSWORD"], config.Db.Password)
	assert.Equal(t, envVars["CSM_DB_NAME"], config.Db.Name)
	assert.Equal(t, envVars["CSM_DB_SSL_MODE"], config.Db.SslMode)
	assert.Equal(t, envVars["CSM_DB_POOL_SIZE"], config.Db.PoolSize)
	assert.Equal(t, envVars["CSM_LOGGING_LEVEL"], config.LoggingLevel)
	assert.Equal(t, envVars["CSM_LOGGING_TYPE"], config.LoggingType)
	assert.Equal(t, envVars["CSM_GRPC_PORT"], config.GrpcServer.Port)
	assert.Equal(t, 1, config.Mcr)
	assert.Equal(t, envVars["CSM_METRICS_PORT"], config.MetricsPort)
	assert.Equal(t, envVars["CSM_PROFILING_PORT"], config.ProfilingPort)

}
