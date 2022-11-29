package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
	"panim.one/nasa/migrations"
	"panim.one/nasa/models"
)

func migrator() {
	dsn := models.BuildPGConnectionString()
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	pg := bun.NewDB(db, pgdialect.New())

	migrator := migrate.NewMigrator(pg, migrations.Migrations)
	err := migrator.Init(context.Background())

	if err != nil {
		fmt.Println(err, "Error on migration init")
	}

	_, err = migrator.Migrate(context.Background())

	if err != nil {
		fmt.Println(err, "Error on migration apply")
	}
}
