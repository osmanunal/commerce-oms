package port

import "context"

type OrderPublisher interface {
	// Stok işlem sonucunu order service'e bildirir
	PublishStockStatus(ctx context.Context, message StockStatusMessage) error
}

type OrderConsumer interface {
	// Order service'den gelen stok düşürme isteklerini dinler
	ConsumeStockDecrease(ctx context.Context, handler func(message StockDecreaseMessage) error) error
}

type OrderBroker interface {
	OrderPublisher
	OrderConsumer
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
