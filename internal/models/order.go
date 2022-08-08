package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// SalesHistory represents a sales history.
type SalesHistory struct {
	gorm.Model
	SalesHistoryItems []SalesHistoryItem
	UserID            string
	TotalPrice        decimal.Decimal
	TotalVat          decimal.Decimal
	TotalDiscount     decimal.Decimal
	SubTotal          decimal.Decimal
}

// SalesHistoryItem represents a sales history item.
type SalesHistoryItem struct {
	gorm.Model
	ProductID      string
	Quantity       int32
	SalesHistoryID uint
}

// NewSalesHistory - creates a new sales history from a shopping cart.
func NewSalesHistory(cart ShoppingCart) SalesHistory {
	return SalesHistory{
		UserID:            cart.UserID,
		TotalPrice:        cart.TotalPrice,
		TotalVat:          cart.TotalVat,
		TotalDiscount:     cart.TotalDiscount,
		SubTotal:          cart.SubTotal,
		SalesHistoryItems: fromShoppingCartItems(cart.Items),
	}
}

// fromShoppingCartItems - converts a shopping cart items to sales history items.
func fromShoppingCartItems(items []ShoppingCartItem) []SalesHistoryItem {
	var orderItems []SalesHistoryItem
	for _, item := range items {
		orderItems = append(orderItems, newSalesHistoryItem(item.ProductID.String(), item.Quantity))
	}
	return orderItems
}

// newSalesHistoryItem - creates a new sales history item with a product ID and quantity.
func newSalesHistoryItem(productID string, quantity int32) SalesHistoryItem {
	return SalesHistoryItem{
		ProductID: productID,
		Quantity:  quantity,
	}
}
