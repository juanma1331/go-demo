package main

import (
	"context"
	"fmt"
	"time"

	auth_domain "github.com/juanma1331/go-demo/internal/auth/domain"
	ecommerce_domain "github.com/juanma1331/go-demo/internal/ecommerce/domain"
	"github.com/juanma1331/go-demo/internal/shared"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
)

func CreateTables(db *bun.DB, ctx context.Context) error {
	_, err := db.NewCreateTable().Model((*auth_domain.User)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*ecommerce_domain.Product)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*ecommerce_domain.Cart)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*ecommerce_domain.CartDetail)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	db, err := shared.OpenDB()

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting migration")
	migrationStart := time.Now()
	createTablesErr := CreateTables(db, context.Background())
	migrationEnd := time.Since(migrationStart)
	fmt.Printf("Migration finished in %s\n", migrationEnd)

	if createTablesErr != nil {
		panic(createTablesErr)
	}
}
