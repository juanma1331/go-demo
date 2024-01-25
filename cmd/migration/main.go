package main

import (
	"context"
	"go-demo/internal/domain"
	"go-demo/internal/infra"

	"github.com/uptrace/bun"
)

func CreateTables(db *bun.DB, ctx context.Context) error {
	_, err := db.NewCreateTable().Model((*domain.User)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*domain.AuthToken)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*domain.Product)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*domain.Cart)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*domain.CartDetail)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := infra.OpenDB(infra.DSN)

	if err != nil {
		panic(err)
	}

	createTablesErr := CreateTables(db, context.Background())

	if createTablesErr != nil {
		panic(createTablesErr)
	}
}
