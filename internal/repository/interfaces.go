package repository

import "user-management-api/internal/models"

type UserRepository interface {
	FindAll()
	FindByUUID()
	Create(user models.User) error
	Update()
	Delete()
	FindByEmail(email string) (models.User, bool)
}
