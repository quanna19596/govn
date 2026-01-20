package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"shopify/internal/config"
	"shopify/internal/db"
	"shopify/internal/db/sqlc"
	"shopify/internal/routes"
	"shopify/internal/validation"
	"shopify/pkg/auth"
	"shopify/pkg/cache"
	"shopify/pkg/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

type ModuleContext struct {
	DB    sqlc.Querier
	Redis *redis.Client
}

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Validator init failed")
	}

	router := gin.Default()

	if err := db.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Database init failed")
	}

	redisClient := config.NewRedisClient()
	cacheRedisService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheRedisService)

	ctx := &ModuleContext{
		DB:    db.DB,
		Redis: redisClient,
	}

	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService, cacheRedisService),
	}

	routes.RegisterRoutes(router, tokenService, cacheRedisService, getModuleRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  router,
		modules: modules,
	}
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
	}

	quit := make(chan os.Signal, 1)
	// syscall.SIGINT --> Ctrl + C
	// syscall.SIGTERM --> Kill service
	// syscall.SIGHUP --> Reload service
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		logger.Log.Info().Msgf("Server is running on %s...", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	<-quit
	logger.Log.Warn().Msg("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Log.Info().Msg("Server exited gracefully")

	return nil
}

func getModuleRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}
