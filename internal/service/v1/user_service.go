package v1service

import (
	"database/sql"
	"errors"
	"shopify/internal/db/sqlc"
	"shopify/internal/repository"
	"shopify/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

func (us *userService) GetAllUsers(ctx *gin.Context, search string, orderBy string, sort string, page int32, limit int32) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()
	if sort == "" {
		sort = "desc"
	}

	if orderBy == "" {
		orderBy = "user_created_at"
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		envLimit := utils.GetEnv("LIMIT_ITEM_ON_PER_PAGE", "10")
		limitInt, err := strconv.Atoi(envLimit)

		if err != nil || limitInt <= 0 {
			limitInt = 10
		}
		limit = int32(limitInt)
	}

	offset := (page - 1) * limit

	users, err := us.repo.GetAllV2(context, search, orderBy, sort, limit, offset)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "failed to fetch users", utils.ErrCodeInternalServer)
	}

	total, err := us.repo.CountUsers(context, search)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "failed to count users", utils.ErrCodeInternalServer)
	}

	return users, int32(total), nil
}

func (us *userService) GetUserByUUID(uuid string) {}

func (us *userService) CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	input.UserEmail = utils.NormalizeString(input.UserEmail)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)

	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternalServer)
	}

	input.UserPassword = string(hashedPassword)

	user, err := us.repo.Create(context, input)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
		}

		return sqlc.User{}, utils.WrapError(err, "failed to create new user", utils.ErrCodeInternalServer)
	}

	return user, nil
}

func (us *userService) UpdateUser(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	if input.UserPassword != nil && *input.UserPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.UserPassword), bcrypt.DefaultCost)

		if err != nil {
			return sqlc.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternalServer)
		}

		hashed := string(hashedPassword)
		input.UserPassword = &hashed
	}

	user, err := us.repo.Update(context, input)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "failed to update user", utils.ErrCodeInternalServer)
	}

	return user, nil
}

func (us *userService) DeleteUser(ctx *gin.Context, uuid uuid.UUID) error {
	context := ctx.Request.Context()

	err := us.repo.Delete(context, uuid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return utils.WrapError(err, "failed to delete user", utils.ErrCodeInternalServer)
	}

	return nil
}

func (us *userService) SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) error {
	context := ctx.Request.Context()

	err := us.repo.SoftDelete(context, uuid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return utils.WrapError(err, "failed to delete user", utils.ErrCodeInternalServer)
	}

	return nil
}

func (us *userService) RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	user, err := us.repo.Restore(context, uuid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "failed to restore user", utils.ErrCodeInternalServer)
	}

	return user, nil
}
