package repository

import (
	"context"
	"fmt"
	"shopify/internal/db"
	"shopify/internal/db/sqlc"

	"github.com/google/uuid"
)

type SqlUserRepository struct {
	DB sqlc.Querier
}

func NewSqlUserRepository(DB sqlc.Querier) UserRepository {
	return &SqlUserRepository{
		DB: DB,
	}
}

func (ur *SqlUserRepository) CountUsers(ctx context.Context, search string, deleted bool) (int64, error) {
	total, err := ur.DB.CountUsers(ctx, sqlc.CountUsersParams{
		Search:  search,
		Deleted: &deleted,
	})

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (ur *SqlUserRepository) GetAll(ctx context.Context, search string, orderBy string, sort string, limit int32, offset int32) ([]sqlc.User, error) {
	var (
		users []sqlc.User
		err   error
	)
	switch {
	case orderBy == "user_id" && sort == "asc":
		users, err = ur.DB.ListUsersUserIdAsc(ctx, sqlc.ListUsersUserIdAscParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})

	case orderBy == "user_id" && sort == "desc":
		users, err = ur.DB.ListUsersUserIdDesc(ctx, sqlc.ListUsersUserIdDescParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})

	case orderBy == "user_created_at" && sort == "asc":
		users, err = ur.DB.ListUsersUserCreatedAtAsc(ctx, sqlc.ListUsersUserCreatedAtAscParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})

	case orderBy == "user_created_at" && sort == "desc":
		users, err = ur.DB.ListUsersUserCreatedAtDesc(ctx, sqlc.ListUsersUserCreatedAtDescParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	}

	if err != nil {
		return []sqlc.User{}, err
	}

	return users, nil
}

func (ur *SqlUserRepository) GetAllV2(ctx context.Context, search string, orderBy string, sort string, limit int32, offset int32, deleted bool) ([]sqlc.User, error) {
	query := `SELECT *
		FROM users
		WHERE (
			$1::TEXT IS NULL
			OR $1::TEXT = ''
			OR user_email ILIKE '%' || $1 || '%'
			OR user_fullname ILIKE '%' || $1 || '%'
		)
	`

	if deleted {
		query += " AND user_deleted_at IS NOT NULL"
	} else {
		query += " AND user_deleted_at IS NULL"
	}

	order := "ASC"
	if sort == "desc" {
		order = "DESC"
	}

	switch orderBy {
	case "user_id", "user_created_at":
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, order)
	default:
		query += " ORDER BY user_id ASC"
	}

	query += " LIMIT $2 OFFSET $3"

	rows, err := db.DBPool.Query(ctx, query, search, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []sqlc.User{}
	for rows.Next() {
		var i sqlc.User
		if err := rows.Scan(
			&i.UserID,
			&i.UserUuid,
			&i.UserEmail,
			&i.UserPassword,
			&i.UserFullname,
			&i.UserAge,
			&i.UserStatus,
			&i.UserLevel,
			&i.UserDeletedAt,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *SqlUserRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.DB.GetUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.DB.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := ur.DB.UpdateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := ur.DB.TrashUser(ctx, uuid)

	if err != nil {
		return err
	}

	return nil
}

func (ur *SqlUserRepository) SoftDelete(ctx context.Context, uuid uuid.UUID) error {
	_, err := ur.DB.SoftDeleteUser(ctx, uuid)

	if err != nil {
		return err
	}

	return nil
}

func (ur *SqlUserRepository) Restore(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.DB.RestoreUser(ctx, uuid)

	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) GetByEmail(ctx context.Context, email string) (sqlc.User, error) {
	user, err := ur.DB.GetUserByEmail(ctx, email)

	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (ur *SqlUserRepository) UpdatePassword(ctx context.Context, input sqlc.UpdatePasswordParams) (sqlc.User, error) {
	user, err := ur.DB.UpdatePassword(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}
