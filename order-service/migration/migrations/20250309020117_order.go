package migrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/osmanunal/commerce-oms/order-service/migration"
	"github.com/osmanunal/commerce-oms/order-service/pkg/model"
	"github.com/shopspring/decimal"

	"github.com/uptrace/bun"
)

type Order struct {
	model.BaseModel

	UserID     string          `bun:",notnull"`
	TotalPrice decimal.Decimal `bun:",notnull"`
	Status     string          `bun:",notnull"`
	OrderItems []OrderItem     `bun:"rel:has-many,join:id=order_id"`
}

type OrderItem struct {
	model.BaseModel

	OrderID   uuid.UUID       `bun:"order_id,type:uuid,notnull"`
	Order     *Order          `bun:"rel:belongs-to,join:order_id=id"`
	ProductID string          `bun:",notnull"`
	Quantity  int             `bun:",notnull"`
	Price     decimal.Decimal `bun:",notnull"`
}

func init() {
	models := []interface{}{
		&Order{},
		&OrderItem{},
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
