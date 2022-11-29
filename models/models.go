package models

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var PG *bun.DB

func Init() {
	dsn := BuildPGConnectionString()

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	PG = bun.NewDB(db, pgdialect.New())
}
