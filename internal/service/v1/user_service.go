package v1service

import (
	"shopify/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUsers(search string, page int, limit int) {}

func (us *userService) GetUserByUUID(uuid string) {}

func (us *userService) CreateUser() {}

func (us *userService) UpdateUser(uuid string) {}

func (us *userService) DeleteUser(uuid string) {}
