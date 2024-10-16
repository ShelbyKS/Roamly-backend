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

	PostgresHost string `mapstructure:"POSTGRES_HOST"`
	PostgresPort string `mapstructure:"POSTGRES_PORT"`
	PostgresUser string `mapstructure:"POSTGRES_USER"`
	PostgresPass string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB   string `mapstructure:"POSTGRES_DB"`
	PostgresSSL  string `mapstructure:"POSTGRES_SSL"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to unmarshal config: %s", err)
	}

	return config
}

func (cfg *Config) GetPostgresCfg() string {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresDB,
		cfg.PostgresPass,
		cfg.PostgresSSL,
	)

	return dsn
}
