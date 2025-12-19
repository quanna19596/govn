package routes

import (
	"shopify/internal/middleware"
	"shopify/internal/utils"
	"shopify/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	httpLogger := newLoggerWithPath("../../internal/logs/http.log", "info")
	recoveryLogger := newLoggerWithPath("../../internal/logs/recovery.log", "warning")
	rateLimiterLogger := newLoggerWithPath("../../internal/logs/rate_limiter.log", "warning")

	r.Use(
		middleware.RateLimiterMiddleware(rateLimiterLogger),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
	)

	v1Api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1Api)
	}
}

func newLoggerWithPath(path string, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:       level,
		Filename:    path,
		MaxSize:     1, // megabytes
		MaxBackups:  5,
		MaxAge:      5, //days
		Compress:    true,
		Environment: utils.GetEnv("APP_ENVIRONMENT", "development"),
	}

	return logger.NewLogger(config)
}
