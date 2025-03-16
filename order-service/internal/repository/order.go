package repository

import (
	"context"
	"errors"
	"github.com/osmanunal/commerce-oms/order-service/internal/model"
	"github.com/osmanunal/commerce-oms/order-service/internal/port"
	"github.com/rs/zerolog/log"

	"github.com/uptrace/bun"
)

type OrderRepository struct {
	db bun.IDB
}

func NewOrderRepository(db *bun.DB) port.OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(order).
			Exec(ctx)

		if err != nil {
			log.Err(err)
			return errors.New("siparişinin oluşturulması sırasında bir hata oluştu")
		}

		for i := range order.OrderItems {
			order.OrderItems[i].OrderID = order.ID
		}

		_, err = tx.NewInsert().
			Model(&order.OrderItems).
			Exec(ctx)

		if err != nil {
			return errors.New("siparişinin oluşturulması sırasında bir hata oluştu")
		}

		return nil
	})
	return err
}

func (r *OrderRepository) Get(ctx context.Context, orderID string) (*model.Order, error) {
	var order model.Order

	err := r.db.NewSelect().
		Model(&order).
		Where("id = ?", orderID).
		Scan(ctx)

	if err != nil {
		return nil, errors.New("siparişinin alınması sırasında bir hata oluştu")
	}

	return &order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	if order == nil {
		return errors.New("siparişinin güncellenmesi sırasında bir hata oluştu")
	}

	_, err := r.db.NewUpdate().
		Model(order).
		Where("id = ?", order.ID).
		Exec(ctx)

	if err != nil {
		return errors.New("siparişinin güncellenmesi sırasında bir hata oluştu")
	}

	return nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID string, status model.OrderStatus) error {
	_, err := r.db.NewUpdate().
		Model(&model.Order{}).
		Where("id = ?", orderID).
		Set("status = ?", status).
		Exec(ctx)

	if err != nil {
		return errors.New("siparişinin durumunun güncellenmesi sırasında bir hata oluştu")
	}

	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, orderID string) error {
	_, err := r.db.NewDelete().
		Model(&model.Order{}).
		Where("id = ?", orderID).
		Exec(ctx)

	if err != nil {
		return errors.New("siparişinin silinmesi sırasında bir hata oluştu")
	}

	return nil
}
