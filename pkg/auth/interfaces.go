package auth

import (
	"shopify/internal/db/sqlc"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	GenerateRefreshToken(user sqlc.User) (RefreshToken, error)
	DecryptAccessTokenPayload(tokenString string) (*EncryptedPayload, error)
	StoreRefreshToken(token RefreshToken) error
	ValidateRefreshToken(token string) (RefreshToken, error)
	RevokeRefreshToken(token string) error
}
