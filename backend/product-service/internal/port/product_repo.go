package port

import (
	"context"

	"github.com/osmanunal/commerce-oms/product-service/internal/model"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]model.Product, error)
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, productID string) (*model.Product, error)
	UpdateStock(ctx context.Context, productID string, stock int) error
}
