package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Secret       string `env:"SECRET"`
	Mail         string `env:"MAIL"`
	KafkaConnect string `env:"KAFKA_PROD" env-default:"kafka:9092"`
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
