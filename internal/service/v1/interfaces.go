package v1service

import (
	"shopify/internal/db/sqlc"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetAllUsers(search string, page int, limit int)
	GetUserByUUID(uuid string)
	CreateUser(ctx *gin.Context, user sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(uuid string)
	DeleteUser(uuid string)
}
