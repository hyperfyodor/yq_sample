package app

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hyperfyodor/yq_sample/internal/config/migrator"
	"github.com/hyperfyodor/yq_sample/internal/helpers"
)

type MigratorApp struct {
	migrate *migrate.Migrate
	config  *migrator.Config
}

func MustLoadMigratorApp() *MigratorApp {
	config := migrator.MustLoad()

	connection := helpers.ConnectionString(
		config.Db.Username,
		config.Db.Password,
		config.Db.Host,
		config.Db.Port,
		config.Db.Name,
		config.Db.SslMode,
	)

	m, err := migrate.New(config.SourceURL, connection)

	if err != nil {
		panic(err)
	}

	return &MigratorApp{m, config}
}

func (app *MigratorApp) Up() error {
	if err := app.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return err
	}

	return nil
}

func (app *MigratorApp) Steps(n int) error {
	if err := app.migrate.Steps(n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return err
	}

	return nil
}

func (app *MigratorApp) Close() {
	sourceErr, databaseErr := app.migrate.Close()

	if databaseErr != nil {
		panic(databaseErr)
	}

	if sourceErr != nil {
		panic(sourceErr)
	}
}
