package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	ConsumedDirPath string `envconfig:"CONSUMED_FOLDER_PATH" default:"./internal/app/testdata/consumed/"`
	Server          struct {
		Host           string `envconfig:"SERVER_HOST" default:"localhost"`
		Port           string `envconfig:"PORT" default:"8080"`
		TimeoutSeconds uint   `envconfig:"SERVER_TIMEOUT" default:"30"`
	}
}

func LoadConfig() (*Config, error) {
	var c Config
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := envconfig.Process("file-pubsub-queue", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
