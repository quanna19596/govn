package repository

import (
	"context"
	"shopify/internal/db/sqlc"
)

type SqlUserRepository struct {
	DB sqlc.Querier
}

func NewSqlUserRepository(DB sqlc.Querier) UserRepository {
	return &SqlUserRepository{
		DB: DB,
	}
}

func (ur *SqlUserRepository) FindAll() {}

func (ur *SqlUserRepository) FindByUUID(uuid string) {}

func (ur *SqlUserRepository) Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.DB.CreateUser(ctx, userParams)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) Update(uuid string) {}

func (ur *SqlUserRepository) Delete(uuid string) {}

func (ur *SqlUserRepository) FindByEmail(email string) {}
