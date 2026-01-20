package routes

import (
	"shopify/internal/middleware"
	v1routes "shopify/internal/routes/v1"
	"shopify/internal/utils"
	"shopify/pkg/auth"
	"shopify/pkg/cache"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, authService auth.TokenService, cacheService cache.RedisCacheService, routes ...Route) {
	httpLogger := utils.NewLoggerWithPath("http.log", "info")
	recoveryLogger := utils.NewLoggerWithPath("recovery.log", "warning")
	rateLimiterLogger := utils.NewLoggerWithPath("rate_limiter.log", "warning")

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(
		middleware.RateLimiterMiddleware(rateLimiterLogger),
		middleware.CORSMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
	)

	v1Api := r.Group("/api/v1")

	middleware.InitAuthMiddleware(authService, cacheService)
	protected := v1Api.Group("")
	protected.Use(
		middleware.AuthMiddleware(),
	)

	for _, route := range routes {
		switch route.(type) {
		case *v1routes.AuthRoutes:
			route.Register(v1Api)
		default:
			route.Register(protected)
		}
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"error": "Not found",
			"path":  ctx.Request.URL.Path,
		})
	})
}
