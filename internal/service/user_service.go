package service

import (
	"strings"
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

func (us *userService) GetAllUsers(search string, page int, limit int) ([]models.User, error) {
	users, err := us.repo.FindAll()
	if err != nil {
		return nil, utils.WrapError(err, "Fail to fetch users", utils.ErrCodeInternalServer)
	}

	var filteredUsers []models.User

	if search != "" {
		search = strings.ToLower(search)
		for _, user := range users {
			name := strings.ToLower(user.Name)
			email := strings.ToLower(user.Email)

			isContainsName := strings.Contains(name, search)
			isContainsEmail := strings.Contains(email, search)

			if isContainsName || isContainsEmail {
				filteredUsers = append(filteredUsers, user)
			}

		}
	} else {
		filteredUsers = users
	}

	start := (page - 1) * limit
	if start >= len(filteredUsers) {
		return []models.User{}, nil
	}

	end := start + limit
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}

	return filteredUsers[start:end], nil
}

func (us *userService) GetUserByUUID(uuid string) (models.User, error) {
	user, found := us.repo.FindByUUID(uuid)

	if !found {
		return models.User{}, utils.NewError("User Not Found", utils.ErrCodeNotFound)
	}

	return user, nil
}

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

func (us *userService) UpdateUser(uuid string, user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)

	if u, alreadyExist := us.repo.FindByEmail(user.Email); alreadyExist && u.UUID != uuid {
		return models.User{}, utils.NewError("Email Already Exist", utils.ErrCodeConflict)
	}

	currentUser, found := us.repo.FindByUUID(uuid)

	if !found {
		return models.User{}, utils.NewError("User Not Found", utils.ErrCodeNotFound)
	}

	currentUser.Name = user.Name
	currentUser.Email = user.Email
	currentUser.Age = user.Age
	currentUser.Status = user.Status

	if err := us.repo.Update(uuid, user); err != nil {
		return models.User{}, utils.WrapError(err, "Fail to update user", utils.ErrCodeInternalServer)
	}

	return currentUser, nil
}

func (us *userService) DeleteUser(uuid string) error {
	if err := us.repo.Delete(uuid); err != nil {
		return utils.WrapError(err, "Fail to delete user", utils.ErrCodeInternalServer)
	}

	return nil
}
