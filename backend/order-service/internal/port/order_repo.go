package port

import (
	"context"
	"github.com/osmanunal/commerce-oms/order-service/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	Get(ctx context.Context, orderID string) (*model.Order, error)
	UpdateStatus(ctx context.Context, orderID string, status model.OrderStatus) error
}
