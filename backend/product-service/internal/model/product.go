package model

import (
	"github.com/osmanunal/commerce-oms/product-service/pkg/model"
	"github.com/shopspring/decimal"
)

type Product struct {
	model.BaseModel

	Name  string          `bun:",notnull"`
	Price decimal.Decimal `bun:",notnull"`
	Stock int             `bun:",notnull"`
}
