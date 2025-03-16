package viewmodel

import (
	"github.com/osmanunal/commerce-oms/product-service/internal/model"
	"github.com/shopspring/decimal"
)

type ProductRequest struct {
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price"`
	Stock int             `json:"stock"`
}

func (p ProductRequest) ToModel(m model.Product) model.Product {
	m.Name = p.Name
	m.Price = p.Price
	m.Stock = p.Stock

	return m
}

type ProductResponse struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price"`
	Stock int             `json:"stock"`
}

func (p *ProductResponse) FromModel(product model.Product) {
	p.ID = product.ID.String()
	p.Name = product.Name
	p.Price = product.Price
	p.Stock = product.Stock
}
