package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

// brokers := []string{"194.87.95.28:29092"}
type Config struct {
	Secret       string `env:"SECRET"`
	Mail         string `env:"MAIL"`
	DSN          string `env:"DATABASE_DSN" env-required:"true"`
	KafkaConnect string `env:"KAFKA_PROD"`
	KafkaTopic   string `env:"KAFKA_TOPIC"`
	KafkaGroup   string `env:"GROUPE"`
}

func Loading(path string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		log.Printf("Warning: не удалось прочитать конфигурационный файл: %v", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
