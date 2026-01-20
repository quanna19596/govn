package v1service

import (
	"shopify/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers(ctx *gin.Context, search string, orderBy string, sort string, page int32, limit int32, deleted bool) ([]sqlc.User, int32, error)
	GetUserByUUID(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	DeleteUser(ctx *gin.Context, uuid uuid.UUID) error
	SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) error
	RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
}

type AuthService interface {
	Login(ctx *gin.Context, email string, password string) (string, string, int, error)
	Logout(ctx *gin.Context, refreshTokenString string) error
	RefreshToken(ctx *gin.Context, refreshTokenString string) (string, string, int, error)
	ForgotPassword(ctx *gin.Context, email string) error
	ResetPassword(ctx *gin.Context, token string, newPassword string) error
}
