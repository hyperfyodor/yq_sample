package producer

import (
	"bytes"
	_ "embed"
	"github.com/ilyakaznacheev/cleanenv"
	"io"
	"log"
)

type Config struct {
	Mps           int                `yaml:"mps" env:"PRD_MPS" env-required:"true" env-description:"messages per second"`
	MaxBacklog    int                `yaml:"max_backlog" env:"PRD_MAX_BACKLOG" env-required:"true" env-description:"highest number of unprocessed messages"`
	LoggingLevel  string             `yaml:"logging_level" env:"PRD_LOGGING_LEVEL" env-required:"true" env-description:"level of the logger"`
	LoggingType   string             `yaml:"logging_type" env:"PRD_LOGGING_TYPE" env-required:"true" env-description:"text - text in console, json - json to the file"`
	Db            dbConfig           `yaml:"db"`
	MetricsPort   string             `yaml:"metrics_port" env:"PRD_METRICS_PORT" env-required:"true" env-description:"port where to scrape metrics"`
	GrpcServer    producerGrpcServer `yaml:"grpc"`
	ProfilingPort string             `yaml:"profiling_port" env:"PRD_PROFILING_PORT" env-required:"true" env-description:"port for profiling"`
}

type producerGrpcServer struct {
	Port string `yaml:"port" env:"PRD_GRPC_PORT" env-required:"true" env-description:"consumer grpc server port"`
	Host string `yaml:"host" env:"PRD_GRPC_HOST" env-required:"true" env-description:"consumer grpc server host"`
}

type dbConfig struct {
	Port     string `yaml:"port" env:"PRD_DB_PORT" env-required:"true" env-description:"db port"`
	Host     string `yaml:"host" env:"PRD_DB_HOST" env-required:"true" env-description:"db host"`
	Name     string `yaml:"name" env:"PRD_DB_NAME" env-required:"true" env-description:"db name"`
	Username string `yaml:"username" env:"PRD_DB_USERNAME" env-required:"true" env-description:"db username"`
	Password string `yaml:"password" env:"PRD_DB_PASSWORD" env-required:"true" env-description:"db password"`
	SslMode  string `yaml:"ssl_mode" env:"PRD_DB_SSL_MODE" env-required:"true" env-description:"ssl mode"`
	PoolSize string `yaml:"pool_size" env:"PRD_DB_POOL_SIZE" env-description:"pool size"`
}

//go:embed config.yml
var producerConfig []byte

func MustLoad() *Config {
	config, err := LoadConfig(bytes.NewBuffer(producerConfig))

	if err != nil {
		log.Printf("Error loading producer config: %s", err)
		panic(err)
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
