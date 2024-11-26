package consumer

import (
	"bytes"
	_ "embed"
	"github.com/ilyakaznacheev/cleanenv"
	"io"
	"log"
)

type Config struct {
	Mcr           int                `yaml:"mcr" env:"CSM_MCR" env-required:"true" env-description:"message consumption rate"`
	LoggingLevel  string             `yaml:"logging_level" env:"CSM_LOGGING_LEVEL" env-required:"true" env-description:"level of the logger"`
	LoggingType   string             `yaml:"logging_type" env:"CSM_LOGGING_TYPE" env-required:"true" env-description:"text - text in console, json - json to the file"`
	Db            dbConfig           `yaml:"db"`
	MetricsPort   string             `yaml:"metrics_port" env:"CSM_METRICS_PORT" env-required:"true" env-description:"port where to scrape metrics"`
	GrpcServer    consumerGrpcServer `yaml:"grpc"`
	ProfilingPort string             `yaml:"profiling_port" env:"CSM_PROFILING_PORT" env-required:"true" env-description:"port for profiling"`
}

type consumerGrpcServer struct {
	Port string `yaml:"port" env:"CSM_GRPC_PORT" env-required:"true" env-description:"consumer grpc server port"`
}

type dbConfig struct {
	Port     string `yaml:"port" env:"CSM_DB_PORT" env-required:"true" env-description:"db port"`
	Host     string `yaml:"host" env:"CSM_DB_HOST" env-required:"true" env-description:"db host"`
	Name     string `yaml:"name" env:"CSM_DB_NAME" env-required:"true" env-description:"db name"`
	Username string `yaml:"username" env:"CSM_DB_USERNAME" env-required:"true" env-description:"db username"`
	Password string `yaml:"password" env:"CSM_DB_PASSWORD" env-required:"true" env-description:"db password"`
	SslMode  string `yaml:"ssl_mode" env:"CSM_DB_SSL_MODE" env-required:"true" env-description:"ssl mode"`
	PoolSize string `yaml:"pool_size" env:"CSM_DB_POOL_SIZE" env-description:"pool size"`
}

//go:embed config.yml
var consumerConfig []byte

func MustLoad() *Config {
	config, err := LoadConsumerConfig(bytes.NewBuffer(consumerConfig))

	if err != nil {
		log.Fatalf("Error loading consumer config: %s", err)
	}

	return config
}

func LoadConsumerConfig(r io.Reader) (*Config, error) {
	var config Config

	if err := cleanenv.ParseYAML(r, &config); err != nil {
		return &Config{}, err
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		return &Config{}, err
	}

	return &config, nil
}
