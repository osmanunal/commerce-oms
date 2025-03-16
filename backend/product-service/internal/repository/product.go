package repository

import (
	"context"
	"errors"

	"github.com/osmanunal/commerce-oms/product-service/internal/model"
	"github.com/osmanunal/commerce-oms/product-service/internal/port"

	"github.com/uptrace/bun"
)

type ProductRepository struct {
	db bun.IDB
}

func NewProductRepository(db *bun.DB) port.ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product

	err := r.db.NewSelect().
		Model(&products).
		Scan(ctx)

	if err != nil {
		return nil, errors.New("ürünlerin alınması sırasında bir hata oluştu")
	}

	return products, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(product).
			Exec(ctx)

		if err != nil {
			return errors.New("ürünün oluşturulması sırasında bir hata oluştu")
		}

		_, err = tx.NewInsert().
			Model(&model.Product{}).
			Exec(ctx)

		if err != nil {
			return errors.New("ürünün oluşturulması sırasında bir hata oluştu")
		}

		return nil
	})
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, productID string) (*model.Product, error) {
	var product model.Product

	err := r.db.NewSelect().
		Model(&product).
		Where("id = ?", productID).
		Scan(ctx)

	if err != nil {
		return nil, errors.New("ürünün alınması sırasında bir hata oluştu")
	}

	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *model.Product) error {
	if product == nil {
		return errors.New("ürünün güncellenmesi sırasında bir hata oluştu")
	}

	_, err := r.db.NewUpdate().
		Model(product).
		Where("id = ?", product.ID).
		Exec(ctx)

	if err != nil {
		return errors.New("ürünün güncellenmesi sırasında bir hata oluştu")
	}

	return nil
}

func (r *ProductRepository) UpdateStock(ctx context.Context, productID string, stock int) error {
	_, err := r.db.NewUpdate().
		Model(&model.Product{}).
		Where("id = ?", productID).
		Set("stock = ?", stock).
		Exec(ctx)

	if err != nil {
		return errors.New("ürünün stokunun güncellenmesi sırasında bir hata oluştu")
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, productID string) error {
	_, err := r.db.NewDelete().
		Model(&model.Product{}).
		Where("id = ?", productID).
		Exec(ctx)

	if err != nil {
		return errors.New("ürünün silinmesi sırasında bir hata oluştu")
	}

	return nil
}
