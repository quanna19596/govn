package auth

import "shopify/internal/db/sqlc"

type TokenService interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	GenerateRefreshToken()
}
