package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shopify/internal/config"
	"shopify/internal/db"
	"shopify/internal/db/sqlc"
	"shopify/internal/routes"
	"shopify/internal/validation"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	DB sqlc.Querier
}

func NewApplication(config *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Validator init failed %v", err)
	}
	loadEnv()

	router := gin.Default()

	if err := db.InitDB(); err != nil {
		log.Fatalf("Database init failed: %v", err)
	}

	ctx := &ModuleContext{
		DB: db.DB,
	}

	modules := []Module{
		NewUserModule(ctx),
	}

	routes.RegisterRoutes(router, getModuleRoutes(modules)...)

	return &Application{
		config:  config,
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
		log.Printf("Server is running on %s... \n", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Faild to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")

	return nil
}

func getModuleRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}

func loadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found")
	}
}
