package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort   string `mapstructure:"SERVER_PORT"`
	LogLevel     string `mapstructure:"LOG_LEVEL"`
	GoogleApiKey string `mapstructure:"GOOGLE_API_KEY"`

	Postgres PostgresConfig `mapstructure:",squash"`
	Redis    RedisConfig    `mapstructure:",squash"`
}

type PostgresConfig struct {
	Host string `mapstructure:"POSTGRES_HOST"`
	Port string `mapstructure:"POSTGRES_PORT"`
	User string `mapstructure:"POSTGRES_USER"`
	Pass string `mapstructure:"POSTGRES_PASSWORD"`
	DB   string `mapstructure:"POSTGRES_DB"`
	SSL  string `mapstructure:"POSTGRES_SSL"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to unmarshal config: %s", err)
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
