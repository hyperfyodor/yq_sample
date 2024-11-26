package migrator

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadMigratorConfigFromEnv(t *testing.T) {
	envVars := map[string]string{
		"DB_HOST":              "test db host",
		"DB_PORT":              "test db port",
		"DB_USERNAME":          "test db user",
		"DB_PASSWORD":          "test db password",
		"DB_NAME":              "test db name",
		"DB_SSL_MODE":          "test db ssl mode",
		"MIGRATION_SOURCE_URL": "test source url",
	}

	for k, v := range envVars {
		err := os.Setenv(k, v)
		if err != nil {
			t.Errorf("Failed to set env var: %s", k)
		}
	}

	config := MustLoad()

	assert.Equal(t, envVars["DB_HOST"], config.Db.Host)
	assert.Equal(t, envVars["DB_PORT"], config.Db.Port)
	assert.Equal(t, envVars["DB_USERNAME"], config.Db.Username)
	assert.Equal(t, envVars["DB_PASSWORD"], config.Db.Password)
	assert.Equal(t, envVars["DB_NAME"], config.Db.Name)
	assert.Equal(t, envVars["DB_SSL_MODE"], config.Db.SslMode)
	assert.Equal(t, envVars["MIGRATION_SOURCE_URL"], config.SourceURL)

}
