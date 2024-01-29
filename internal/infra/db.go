package infra

import (
	"database/sql"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func OpenDB() (*bun.DB, error) {
	dsn := os.Getenv("DATABASE_URL")

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := db.Ping(); err != nil {
		return nil, err
	}

	bunDB := bun.NewDB(db, pgdialect.New())
	// bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return bunDB, nil
}
