package infra

import (
	"context"
	"database/sql"

	"github.com/juanma1331/go-demo/internal/auth/app/services"
	"github.com/juanma1331/go-demo/internal/auth/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type userRepository struct {
	*bun.DB
}

func NewUserRepository(db *bun.DB) *userRepository {
	return &userRepository{db}
}

func (bus *userRepository) InsertUserByEmail(user *domain.User) error {
	_, err := bus.NewInsert().Model(user).Exec(context.Background())

	return err
}

func (bus *userRepository) SelectUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := bus.NewSelect().Model(&user).Where("email = ?", email).Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, services.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (bus *userRepository) SelectUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := bus.NewSelect().Model(&user).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, services.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
