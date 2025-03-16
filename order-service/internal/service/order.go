package service

import (
	"context"
	"github.com/osmanunal/commerce-oms/order-service/internal/model"
	"github.com/osmanunal/commerce-oms/order-service/internal/port"
	"log"
)

type OrderService struct {
	orderRepo     port.OrderRepository
	productBroker port.ProductBroker
}

func NewOrderService(orderRepo port.OrderRepository, productBroker port.ProductBroker) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		productBroker: productBroker,
	}
}

func (s *OrderService) Create(ctx context.Context, order *model.Order) error {
	order.Status = model.OrderStatusPending
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return err
	}

	stockDecreaseMessages := make([]port.StockDecreaseMessage, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		stockDecreaseMessages = append(stockDecreaseMessages, port.StockDecreaseMessage{
			OrderID:   order.ID.String(),
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	if err := s.productBroker.PublishStockDecrease(ctx, stockDecreaseMessages); err != nil {
		err = s.orderRepo.UpdateStatus(ctx, order.ID.String(), model.OrderStatusFailed)
		return err
	}

	return nil
}

func (s *OrderService) Get(ctx context.Context, orderID string) (*model.Order, error) {
	return s.orderRepo.Get(ctx, orderID)
}

func (s *OrderService) StartConsuming(ctx context.Context) error {
	return s.productBroker.ConsumeStockStatus(ctx, func(message port.StockStatusMessage) error {
		var status model.OrderStatus
		if message.Success {
			status = model.OrderStatusCreated
		} else {
			log.Printf("stok işlemi başarısız: %s - hata: %s", message.OrderID, message.ErrorMsg)
			status = model.OrderStatusFailed
		}

		return s.orderRepo.UpdateStatus(ctx, message.OrderID, status)
	})
}
