// Package config предоставляет структуру и загрузку конфигурации из переменных окружения.
package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config содержит все параметры конфигурации, загружаемые из переменных окружения.
type Config struct {
	Env             string        `env:"ENV" envDefault:"dev"`
	LogsPath        string        `env:"LOGS_PATH"`
	Version         string        `env:"VERSION"`
	HTTPPort        string        `env:"HTTP_PORT" envDefault:":8080"`
	HTTPTimeout     time.Duration `env:"HTTP_TIMEOUT" envDefault:"4s"`
	HTTPIdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`

	PostgresURL     string `env:"POSTGRES_URL"`
	PostgresHost    string `env:"POSTGRES_HOST"`
	PostgresPort    uint16 `env:"POSTGRES_PORT"`
	PostgresUser    string `env:"POSTGRES_USER"`
	PostgresPass    string `env:"POSTGRES_PASSWORD"`
	PostgresDB      string `env:"POSTGRES_DB"`
	PostgresPoolMax int    `env:"POSTGRES_POOL_MAX" envDefault:"5"`

	CacheCap int `env:"CACHE_CAPACITY" envDefault:"10"`

	KafkaBrokers     []string      `env:"KAFKA_BROKERS" envSeparator:","`
	KafkaTopic       string        `env:"KAFKA_TOPIC" envDefault:"orders"`
	KafkaGroupID     string        `env:"KAFKA_GROUP_ID" envDefault:"wb_orders_consumer"`
	KafkaMinBytes    int           `env:"KAFKA_MIN_BYTES" envDefault:"1"`
	KafkaMaxBytes    int           `env:"KAFKA_MAX_BYTES" envDefault:"10485760"`
	KafkaMaxWait     time.Duration `env:"KAFKA_MAX_WAIT" envDefault:"1s"`
	KafkaReadTimeout time.Duration `env:"KAFKA_READ_TIMEOUT" envDefault:"5s"`
	KafkaDialTimeout time.Duration `env:"KAFKA_DIAL_TIMEOUT" envDefault:"5s"`
	KafkaMsgTimeout  time.Duration `env:"KAFKA_MSG_TIMEOUT" envDefault:"3s"`
}

// MustLoad парсит переменные окружения и возвращает конфигурацию или завершает выполнение при ошибке.
func MustLoad() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Cant load configuration: %v", err)
	}

	return cfg
}
