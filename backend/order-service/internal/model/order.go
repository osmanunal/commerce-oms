package model

import (
	"github.com/google/uuid"
	"github.com/osmanunal/commerce-oms/order-service/pkg/model"
	"github.com/shopspring/decimal"
)

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "pending"
	OrderStatusCreated OrderStatus = "created"
	OrderStatusFailed  OrderStatus = "failed"
)

type Order struct {
	model.BaseModel

	UserID     string          `bun:",notnull"`
	TotalPrice decimal.Decimal `bun:",notnull"`
	Status     OrderStatus     `bun:",notnull"`
	OrderItems []OrderItem     `bun:"rel:has-many,join:id=order_id"`
}

type OrderItem struct {
	model.BaseModel

	OrderID   uuid.UUID       `bun:"order_id,notnull"`
	Order     *Order          `bun:"rel:belongs-to,join:order_id=id"`
	ProductID string          `bun:",notnull"`
	Quantity  int             `bun:",notnull"`
	Price     decimal.Decimal `bun:",notnull"`
}
