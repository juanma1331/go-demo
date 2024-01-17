package repositories

import (
	"context"
	"database/sql"
	"errors"
	"go-demo/internal/app/services"
	"go-demo/internal/domain"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type sqliteAuthTokenRepository struct {
	*bun.DB
}

func NewSqliteAuthTokenRepository(db *bun.DB) *sqliteAuthTokenRepository {
	return &sqliteAuthTokenRepository{db}
}

func (atr *sqliteAuthTokenRepository) InsertToken(t *domain.AuthToken) error {
	t.ID = uuid.New()
	_, err := atr.NewInsert().Model(t).Exec(context.Background())

	return err
}

func (atr *sqliteAuthTokenRepository) SelectToken(token string) (domain.AuthToken, error) {
	var at domain.AuthToken

	err := atr.NewSelect().Model(&at).Where("token = ?", token).Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return at, services.ErrTokenNotFound
		}

		return at, err
	}

	return at, nil
}

func (atr *sqliteAuthTokenRepository) DeleteToken(token string) error {
	_, err := atr.NewDelete().Model(&domain.AuthToken{}).Where("token = ?", token).Exec(context.Background())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return services.ErrTokenNotFound
		}

	}

	return err
}
