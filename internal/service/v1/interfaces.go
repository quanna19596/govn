package v1service

import (
	"shopify/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers(ctx *gin.Context, search string, orderBy string, sort string, page int32, limit int32) ([]sqlc.User, int32, error)
	GetUserByUUID(uuid string)
	CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	DeleteUser(ctx *gin.Context, uuid uuid.UUID) error
	SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) error
	RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
}
