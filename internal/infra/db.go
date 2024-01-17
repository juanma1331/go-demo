package infra

import (
	"database/sql"

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

	bunDB := bun.NewDB(db, sqlitedialect.New())

	if err != nil {
		return nil, err
	}

	return bunDB, nil
}
