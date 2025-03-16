package service

import (
	"context"
	"fmt"
	"log"

	"github.com/osmanunal/commerce-oms/product-service/internal/model"
	"github.com/osmanunal/commerce-oms/product-service/internal/port"
)

type ProductService struct {
	productRepo port.ProductRepository
	orderBroker port.OrderBroker
}

func NewProductService(productRepo port.ProductRepository, orderBroker port.OrderBroker) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		orderBroker: orderBroker,
	}
}

func (s *ProductService) Create(ctx context.Context, product *model.Product) error {
	return s.productRepo.Create(ctx, product)
}

func (s *ProductService) Get(ctx context.Context, productID string) (*model.Product, error) {
	return s.productRepo.Get(ctx, productID)
}

func (s *ProductService) StartConsuming(ctx context.Context) error {
	return s.orderBroker.ConsumeStockDecrease(ctx, func(message port.StockDecreaseMessage) error {
		product, err := s.productRepo.Get(ctx, message.ProductID)
		if err != nil {
			return err
		}

		if product.Stock < message.Quantity {
			return fmt.Errorf("yetersiz stok: mevcut %d, istenen %d", product.Stock, message.Quantity)
		}

		newStock := product.Stock - message.Quantity
		err = s.productRepo.UpdateStock(ctx, message.ProductID, newStock)

		statusMessage := port.StockStatusMessage{
			OrderID: message.OrderID,
			Success: err == nil,
		}

		if err != nil {
			statusMessage.ErrorMsg = err.Error()
			log.Printf("stok düşürme hatası: %v", err)
		}

		return s.orderBroker.PublishStockStatus(ctx, statusMessage)
	})
}
