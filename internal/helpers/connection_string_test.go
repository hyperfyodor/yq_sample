package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnectionString(t *testing.T) {
	username := "username"
	password := "password"
	host := "host"
	port := "port"
	name := "name"
	sslMode := "sslmode"

	connectionString := ConnectionString(username, password, host, port, name, sslMode)

	assert.Equal(t, "postgres://username:password@host:port/name?sslmode=sslmode", connectionString)
}
