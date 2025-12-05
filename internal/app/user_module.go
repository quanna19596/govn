package app

import (
	"user-management-api/internal/handler"
	"user-management-api/internal/repository"
	"user-management-api/internal/routes"
	"user-management-api/internal/service"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	repo := repository.NewInMemoryUserRepository()
	service := service.NewUserService(repo)
	handler := handler.NewUserHandler(service)
	routes := routes.NewUserRoutes(handler)

	return &UserModule{routes: routes}
}

func (um *UserModule) Routes() routes.Route {
	return um.routes
}
