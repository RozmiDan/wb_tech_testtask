package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
)

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
	// Kafka
	// KafkaBrokers []string      `env:"KAFKA_BROKERS" envSeparator:","`
	// KafkaTopic   string        `env:"KAFKA_TOPIC"`
	// KafkaTimeout time.Duration `env:"KAFKA_TIMEOUT" envDefault:"5s"`
}

func MustLoad() *Config {
	var cfg *Config = &Config{}

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Cant load configuration: %v", err)
	}
	return cfg
}
