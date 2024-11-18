package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerPort   string `envconfig:"SERVER_PORT"`
	LogLevel     string `envconfig:"LOG_LEVEL"`
	GoogleApiKey string `envconfig:"GOOGLE_API_KEY"`
	OpenAiKey    string `envconfig:"OPEN_AI_KEY"`
	JWTSecret    string `envconfig:"JWT_SECRET"`

	Postgres PostgresConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	Host string `envconfig:"POSTGRES_HOST"`
	Port string `envconfig:"POSTGRES_PORT"`
	User string `envconfig:"POSTGRES_USER"`
	Pass string `envconfig:"POSTGRES_PASSWORD"`
	DB   string `envconfig:"POSTGRES_DB"`
	SSL  string `envconfig:"POSTGRES_SSL"`
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST"`
	Port     string `envconfig:"REDIS_PORT"`
	Password string `envconfig:"REDIS_PASSWORD"`
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

func (cfg *Config) GetPostgresCfg() string {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.DB,
		cfg.Postgres.Pass,
		cfg.Postgres.SSL,
	)

	return dsn
}
