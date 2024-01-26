package infra

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

const DSN = "file:data/data.sqlite"

func OpenDB(dsn string) (*bun.DB, error) {
	db, err := sql.Open(sqliteshim.ShimName, dsn)

	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Set WAL journal mode
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, err
	}

	bunDB := bun.NewDB(db, sqlitedialect.New())

	if err != nil {
		return nil, err
	}

	return bunDB, nil
}
