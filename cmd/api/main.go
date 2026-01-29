package main

import (
	"path/filepath"
	"shopify/internal/app"
	"shopify/internal/config"
	"shopify/internal/utils"
	"shopify/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	rootDir := utils.MustGetWorkingDir()

	logFile := filepath.Join(rootDir, "internal/logs/app.log")

	logger.InitLogger(logger.LoggerConfig{
		Level:       "info",
		Filename:    logFile,
		MaxSize:     1,
		MaxBackups:  5,
		MaxAge:      5,
		Compress:    true,
		Environment: utils.GetEnv("APP_ENVIRONMENT", "development"),
	})

	loadEnv(filepath.Join(rootDir, ".env"))

	config := config.NewConfig()
	application, err := app.NewApplication(config)

	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize application")
	}

	if err := application.Run(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Application run failed")
	}
}

func loadEnv(path string) {
	if err := godotenv.Load(path); err != nil {
		logger.Log.Warn().Msg("⚠️ No .env file found")
	} else {
		logger.Log.Info().Msg("✅ Loaded successfully .env in api proccess")
	}
}
