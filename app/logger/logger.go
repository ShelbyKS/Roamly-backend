package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ShelbyKS/Roamly-backend/app/config"
)

func InitLogger(config *config.Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{})

	logLevel, err := logrus.ParseLevel(strings.ToLower(config.LogLevel))
	if err != nil {
		logger.Warn("Invalid log level '%s', defaulting to 'info'", config.LogLevel)
		logLevel = logrus.InfoLevel
	}

	logger.SetLevel(logLevel)

	return logger
}
