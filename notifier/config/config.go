package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	ServerPort  string `envconfig:"NOTIFIER_PORT"`
	KafkaConfig KafkaConfig
}

type KafkaConfig struct {
	Host  string `envconfig:"KAFKA_HOST"`
	Port  string `envconfig:"KAFKA_PORT"`
	Topic string `envconfig:"KAFKA_TOPIC"`
	Group string `envconfig:"KAFKA_GROUP"`
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found, proceeding with environment variables")
	}

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Unable to process environment variables: %s", err)
	}

	return &config
}
