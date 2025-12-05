package service

import (
	"user-management-api/internal/models"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUsers() {}

func (us *userService) GetUserByUUID() {}

func (us *userService) CreateUser(user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)

	if _, alreadyExist := us.repo.FindByEmail(user.Email); alreadyExist {
		return models.User{}, utils.NewError("Email Already Exist", utils.ErrCodeConflict)
	}

	user.UUID = uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return models.User{}, utils.WrapError(err, "Fail to hash password", utils.ErrCodeInternalServer)
	}

	user.Password = string(hashedPassword)

	if err := us.repo.Create(user); err != nil {
		return models.User{}, utils.WrapError(err, "Fail to create user", utils.ErrCodeInternalServer)
	}

	return user, nil
}

func (us *userService) UpdateUser() {}

func (us *userService) DeleteUser() {}
