package service

import (
	"log"
	"user-management-api/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo,
	}
}

func (us *userService) GetAllUsers() {
	log.Println("GetAllUsers service")
	us.repo.FindAll()
}

func (us *userService) GetUserByUUID() {}

func (us *userService) CreateUser() {}

func (us *userService) UpdateUser() {}

func (us *userService) DeleteUser() {}
