package app

import (
	v1handler "shopify/internal/handler/v1"
	"shopify/internal/repository"
	"shopify/internal/routes"
	v1routes "shopify/internal/routes/v1"
	v1service "shopify/internal/service/v1"
	"shopify/pkg/auth"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(ctx *ModuleContext, tokenService auth.TokenService) *AuthModule {
	repo := repository.NewSqlUserRepository(ctx.DB)
	service := v1service.NewAuthService(repo, tokenService)
	handler := v1handler.NewAuthHandler(service)
	routes := v1routes.NewAuthRoutes(handler)

	return &AuthModule{routes: routes}
}

func (am *AuthModule) Routes() routes.Route {
	return am.routes
}
