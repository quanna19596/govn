package service

import "user-management-api/internal/models"

type UserService interface {
	GetAllUsers()
	GetUserByUUID()
	CreateUser(user models.User) (models.User, error)
	UpdateUser()
	DeleteUser()
}
