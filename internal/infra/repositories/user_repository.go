package repositories

import (
	"context"
	"database/sql"
	"go-demo/internal/app/services"
	"go-demo/internal/domain"

	"github.com/uptrace/bun"
)

type sqliteUserRepository struct {
	*bun.DB
}

func NewSqliteUserRepository(db *bun.DB) *sqliteUserRepository {
	return &sqliteUserRepository{db}
}

func (bus *sqliteUserRepository) InsertUserByEmail(user *domain.User) error {
	_, err := bus.NewInsert().Model(user).Exec(context.Background())

	return err
}

func (bus *sqliteUserRepository) SelectUserByEmail(email string) (*domain.User, error) {
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

func (bus *sqliteUserRepository) SelectUserByID(id string) (*domain.User, error) {
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
