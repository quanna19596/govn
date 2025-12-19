package repository

import (
	"context"
	"shopify/internal/db/sqlc"
)

type UserRepository interface {
	FindAll()
	FindByUUID(uuid string)
	Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	Update(uuid string)
	Delete(uuid string)
	FindByEmail(email string)
}
