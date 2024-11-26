package migrator

import (
	"bytes"
	_ "embed"
	"github.com/ilyakaznacheev/cleanenv"
	"io"
	"log"
)

type Config struct {
	Db        dbConfig `yaml:"db"`
	SourceURL string   `yaml:"source_url" env:"MIGRATION_SOURCE_URL" env-required:"true" env-description:"Source URL for migrator"`
}

type dbConfig struct {
	Port     string `yaml:"port" env:"DB_PORT" env-required:"true" env-description:"db port"`
	Host     string `yaml:"host" env:"DB_HOST" env-required:"true" env-description:"db host"`
	Name     string `yaml:"name" env:"DB_NAME" env-required:"true" env-description:"db name"`
	Username string `yaml:"username" env:"DB_USERNAME" env-required:"true" env-description:"db username"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true" env-description:"db password"`
	SslMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE" env-required:"true" env-description:"ssl mode"`
}

//go:embed config.yml
var migratorConfig []byte

func MustLoad() *Config {
	config, err := LoadConfig(bytes.NewBuffer(migratorConfig))

	if err != nil {
		log.Fatalf("Error loading migrator config: %s", err)
	}

	return config
}

func LoadConfig(r io.Reader) (*Config, error) {
	var config Config

	if err := cleanenv.ParseYAML(r, &config); err != nil {
		return &Config{}, err
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		return &Config{}, err
	}

	return &config, nil
}
