package models

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// ShoppingCartItem - represents a shopping cart item.
type ShoppingCartItem struct {
	Base
	ProductID      uuid.UUID       `json:"product_id"`
	ProductName    string          `json:"product_name"`
	Quantity       int32           `json:"quantity"`
	Price          decimal.Decimal `json:"price"`
	VatRate        int32           `json:"vat_rate"`
	ShoppingCartID string          `json:"shopping_cart_id"`
}

// NewShoppingCartItem - creates a new shopping cart item from a product ID and quantity and shopping cart ID and vat rate.
func NewShoppingCartItem(productID uuid.UUID, productName string, quantity int32, price decimal.Decimal, vatRate int32, shoppingCartID string) ShoppingCartItem {
	cartItemId := uuid.Must(uuid.NewV4())
	return ShoppingCartItem{
		Base: Base{
			ID: cartItemId,
		},
		ProductID:      productID,
		ProductName:    productName,
		Quantity:       quantity,
		Price:          price,
		VatRate:        vatRate,
		ShoppingCartID: shoppingCartID,
	}
}
