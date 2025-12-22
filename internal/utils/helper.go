package utils

import (
	"os"
	"shopify/pkg/logger"
	"strconv"

	"github.com/rs/zerolog"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intVal
}

func NewLoggerWithPath(path string, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:       level,
		Filename:    path,
		MaxSize:     1, // megabytes
		MaxBackups:  5,
		MaxAge:      5, //days
		Compress:    true,
		Environment: GetEnv("APP_ENVIRONMENT", "development"),
	}

	return logger.NewLogger(config)
}
