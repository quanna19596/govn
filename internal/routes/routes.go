package routes

import (
	"shopify/internal/middleware"
	"shopify/internal/utils"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	httpLogger := utils.NewLoggerWithPath("../../internal/logs/http.log", "info")
	recoveryLogger := utils.NewLoggerWithPath("../../internal/logs/recovery.log", "warning")
	rateLimiterLogger := utils.NewLoggerWithPath("../../internal/logs/rate_limiter.log", "warning")

	r.Use(
		middleware.RateLimiterMiddleware(rateLimiterLogger),
		middleware.CORSMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
	)

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	v1Api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1Api)
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"error": "Not found",
			"path":  ctx.Request.URL.Path,
		})
	})
}
