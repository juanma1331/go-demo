package repositories

import (
	"context"
	"database/sql"
	"go-demo/internal/app/services/authservice"
	"go-demo/internal/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type sqliteAuthUserRepository struct {
	*bun.DB
}

func NewSqliteUserRepository(db *bun.DB) *sqliteAuthUserRepository {
	return &sqliteAuthUserRepository{db}
}

func (bus *sqliteAuthUserRepository) InsertUserByEmail(user *domain.User) error {
	_, err := bus.NewInsert().Model(user).Exec(context.Background())

	return err
}

func (bus *sqliteAuthUserRepository) SelectUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := bus.NewSelect().Model(&user).Where("email = ?", email).Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, authservice.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (bus *sqliteAuthUserRepository) SelectUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := bus.NewSelect().Model(&user).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, authservice.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
