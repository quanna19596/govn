package repository

import (
	"context"
	"shopify/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	CountUsers(ctx context.Context, search string, deleted bool) (int64, error)
	GetAll(ctx context.Context, search string, orderBy string, sort string, limit int32, offset int32) ([]sqlc.User, error)
	GetAllV2(ctx context.Context, search string, orderBy string, sort string, limit int32, offset int32, deleted bool) ([]sqlc.User, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	Update(ctx context.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	SoftDelete(ctx context.Context, uuid uuid.UUID) error
	Restore(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	GetByEmail(ctx context.Context, email string) (sqlc.User, error)
}
