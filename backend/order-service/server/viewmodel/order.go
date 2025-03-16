package viewmodel

import (
	"github.com/osmanunal/commerce-oms/order-service/internal/model"
	"github.com/shopspring/decimal"
)

type OrderRequest struct {
	UserID     string             `json:"user_id"` //Bu alan normalde jwt token'dan alınacak
	OrderItems []OrderItemRequest `json:"order_items" validate:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID string          `json:"product_id" validate:"required"`
	Quantity  int             `json:"quantity" validate:"required,min=1"`
	Price     decimal.Decimal `json:"price"` // Bu alan normalde product service'ten alınacak
}

func (o OrderRequest) ToModel(m model.Order) model.Order {
	var orderItems []model.OrderItem
	var totalPrice decimal.Decimal
	for _, item := range o.OrderItems {
		orderItems = append(orderItems, model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})

		itemTotal := item.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))
		totalPrice = totalPrice.Add(itemTotal)
	}

	m.UserID = o.UserID
	m.OrderItems = orderItems
	m.TotalPrice = totalPrice
	return m
}

type OrderResponse struct {
	ID         string              `json:"id"`
	UserID     string              `json:"user_id"`
	OrderItems []OrderItemResponse `json:"order_items"`
	TotalPrice decimal.Decimal     `json:"total_price"`
	Status     string              `json:"status"`
}

type OrderItemResponse struct {
	ProductID string          `json:"product_id"`
	Quantity  int             `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
	Status    string          `json:"status"`
}

func (o OrderResponse) FromModel(order model.Order) OrderResponse {
	return OrderResponse{
		ID:         order.ID.String(),
		UserID:     order.UserID,
		OrderItems: o.OrderItems,
		TotalPrice: order.TotalPrice,
		Status:     string(order.Status),
	}
}
