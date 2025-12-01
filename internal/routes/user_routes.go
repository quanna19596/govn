package routes

import (
	"user-management-api/internal/handler"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *handler.UserHandler
}

func NewUserRoutes(handler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler,
	}
}

func (ur *UserRoutes) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("", ur.handler.GetAllUsers)
		users.GET("/:uuid", ur.handler.GetUserByUUID)
		users.POST("", ur.handler.CreateUser)
		users.PUT("", ur.handler.UpdateUser)
		users.DELETE("", ur.handler.DeleteUser)

	}
}
