package port

import (
	"context"
)

type ProductBroker interface {
	PublishStockDecrease(ctx context.Context, message []StockDecreaseMessage) error
	ConsumeStockStatus(ctx context.Context, handler func(message StockStatusMessage) error) error
}

type StockDecreaseMessage struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type StockStatusMessage struct {
	OrderID  string `json:"order_id"`
	Success  bool   `json:"success"`
	ErrorMsg string `json:"error_msg,omitempty"`
}
