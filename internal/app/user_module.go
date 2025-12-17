package app

import (
	v1handler "shopify/internal/handler/v1"
	"shopify/internal/repository"
	"shopify/internal/routes"
	v1routes "shopify/internal/routes/v1"
	v1service "shopify/internal/service/v1"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	repo := repository.NewSqlUserRepository()
	service := v1service.NewUserService(repo)
	handler := v1handler.NewUserHandler(service)
	routes := v1routes.NewUserRoutes(handler)

	return &UserModule{routes: routes}
}

func (um *UserModule) Routes() routes.Route {
	return um.routes
}
