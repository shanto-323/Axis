package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shanto-323/axis/internal/model/dto"
	"github.com/shanto-323/axis/internal/model/entity"
)

func (db *DB) CreateUser(ctx context.Context, userDto *dto.RegisterRequest) (*entity.User, error) {
	userEntity := entity.User{
		Email:        userDto.Email,
		PasswordHash: userDto.Password,
	}

	userEntity.ID = uuid.New()
	userEntity.CreatedAt = time.Now()
	userEntity.UpdatedAt = time.Now()

	db.pool[userDto.Email] = &userEntity
	db.logger.Info().
		Str("event", "user_created").
		Str("email", userDto.Email).
		Msg("new user created")

	return &userEntity, nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, ok := db.pool[email]
	if !ok {
		db.logger.Warn().Msg("no user found")
		return nil, fmt.Errorf("no user found")
	}

	userType, _ := user.(*entity.User)

	return userType, nil
}
