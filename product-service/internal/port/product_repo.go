package port

import (
	"context"

	"github.com/osmanunal/commerce-oms/product-service/internal/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	Get(ctx context.Context, productID string) (*model.Product, error)
	UpdateStock(ctx context.Context, productID string, stock int) error
}
