package storage

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hyperfyodor/yq_sample/internal/app/migrator"
	"github.com/hyperfyodor/yq_sample/internal/storage"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
)

type Suite struct {
	*testing.T
	Storage *storage.PostgresStorage
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Minute)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	container := StartDockerPostgres(ctx, "integration", "integration", "tasks")
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	envVars := map[string]string{
		"DB_HOST":              "localhost",
		"DB_PORT":              port.Port(),
		"DB_USERNAME":          "integration",
		"DB_PASSWORD":          "integration",
		"DB_NAME":              "tasks",
		"DB_SSL_MODE":          "disable",
		"MIGRATION_SOURCE_URL": "file://../migration",
	}

	setupEnv(envVars)

	migratorApp := app.MustLoadMigratorApp()
	err = migratorApp.Up()

	if err != nil {
		t.Logf("failed to migrate: %v", err)
		migratorApp.Close()
		panic(err)
	}

	migratorApp.Close()

	connection := container.MustConnectionString(ctx) + "sslmode=disable" + "&pool_max_conns=1"

	pg, err := storage.NewPostgresStorage(ctx, connection, true)

	if err != nil {
		t.Logf("failed to create storage: %v", err)
		panic(err)
	}

	return ctx, &Suite{Storage: pg}

}

func StartDockerPostgres(ctx context.Context, username, password, databaseName string) *postgres.PostgresContainer {
	c, err := postgres.Run(ctx,
		"postgres:17.1-alpine3.20",
		postgres.WithDatabase(databaseName),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		panic(err)
	}

	return c
}

func setupEnv(kv map[string]string) {
	for k, v := range kv {
		err := os.Setenv(k, v)
		if err != nil {
			panic(fmt.Errorf("failed to set env var: %s", k))
		}
	}
}
