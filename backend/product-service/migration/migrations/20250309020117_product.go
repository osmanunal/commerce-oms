package migrations

import (
	"context"

	"github.com/osmanunal/commerce-oms/product-service/migration"
	"github.com/osmanunal/commerce-oms/product-service/pkg/model"
	"github.com/shopspring/decimal"

	"github.com/uptrace/bun"
)

type Product struct {
	model.BaseModel

	Name  string          `bun:",notnull"`
	Price decimal.Decimal `bun:",notnull"`
	Stock int             `bun:",notnull"`
}

func init() {
	models := []interface{}{
		&Product{},
	}
	up := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
			if err != nil {
				return err
			}

			for _, m := range models {
				_, err = tx.NewCreateTable().Model(m).IfNotExists().WithForeignKeys().Exec(ctx)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	down := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) (err error) {
			for _, m := range models {
				_, err = tx.NewDropTable().Model(m).IfExists().Cascade().Exec(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	migration.Migrations.MustRegister(up, down)
}
