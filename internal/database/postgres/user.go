package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/shanto-323/axis/internal/errs"
	"github.com/shanto-323/axis/internal/model/dto"
	"github.com/shanto-323/axis/internal/model/entity"
)

func (db *DB) CreateUser(ctx context.Context, userDto *dto.RegisterRequest) (*entity.User, error) {
	query := `
		INSERT INTO users (
			email,
			password
		)
		VALUES (
			@email,
			@password
		)
		RETURNING 
			id,
			email,
			password,
			created_at,
			updated_at
	`

	user := &entity.User{}

	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"email":    userDto.Email,
		"password": userDto.Password,
	}).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			code := "USER_NOT_FOUND"
			return nil, errs.NewNotFoundError(err.Error(), true, &code)
		}
		return nil, err
	}

	return user, nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT
			id,
			password
		FROM 
			users
		WHERE 
			email = @email
	`

	// only password
	user := &entity.User{}

	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{"email": email}).Scan(
		&user.ID,
		&user.PasswordHash,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			code := "USER_NOT_FOUND"
			return nil, errs.NewNotFoundError(err.Error(), true, &code)
		}
		return nil, err
	}

	return user, nil
}
