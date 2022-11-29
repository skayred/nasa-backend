package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"panim.one/nasa/models"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().
			Model((*models.Asteroid)(nil)).
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		_, err = db.NewCreateTable().
			Model((*models.ProximityEvent)(nil)).
			WithForeignKeys().
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
