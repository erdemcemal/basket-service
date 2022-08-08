package dto

import "github.com/shopspring/decimal"

type ProductDTO struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	UnitPrice decimal.Decimal `json:"price"`
	VatRate   int32           `json:"vatRate"`
	Quantity  int32           `json:"quantity"`
}

type ShoppingCartDTO struct {
	ID            string                `json:"id"`
	UserID        string                `json:"user_id"`
	Items         []ShoppingCartItemDTO `json:"items"`
	TotalPrice    decimal.Decimal       `json:"total_price"`
	TotalVat      decimal.Decimal       `json:"total_vat"`
	TotalDiscount decimal.Decimal       `json:"total_discount"`
	SubTotal      decimal.Decimal       `json:"sub_total"`
}

type ShoppingCartItemDTO struct {
	ID        string          `json:"id"`
	ProductID string          `json:"product_id"`
	Quantity  int32           `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
	VatRate   int32           `json:"vat_rate"`
	Name      string          `json:"name"`
}

type AddItemToBasketDTO struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int32  `json:"quantity" validate:"gte=1,required"`
}

type UpdateItemInBasketDTO struct {
	Quantity  int32  `json:"quantity" validate:"gte=1,required"`
	ProductID string `json:"product_id" validate:"required"`
}
