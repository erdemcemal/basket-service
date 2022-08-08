package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type SalesHistory struct {
	gorm.Model
	SalesHistoryItems []SalesHistoryItem
	UserID            string
	TotalPrice        decimal.Decimal
	TotalVat          decimal.Decimal
	TotalDiscount     decimal.Decimal
	SubTotal          decimal.Decimal
}

type SalesHistoryItem struct {
	gorm.Model
	ProductID      string
	Quantity       int32
	SalesHistoryID uint
}

func NewOrderHistory(cart ShoppingCart) SalesHistory {
	return SalesHistory{
		UserID:            cart.UserID,
		TotalPrice:        cart.TotalPrice,
		TotalVat:          cart.TotalVat,
		TotalDiscount:     cart.TotalDiscount,
		SubTotal:          cart.SubTotal,
		SalesHistoryItems: fromShoppingCartItems(cart.Items),
	}
}

func fromShoppingCartItems(items []ShoppingCartItem) []SalesHistoryItem {
	var orderItems []SalesHistoryItem
	for _, item := range items {
		orderItems = append(orderItems, newOrderHistoryItem(item.ProductID.String(), item.Quantity))
	}
	return orderItems
}

func newOrderHistoryItem(productID string, quantity int32) SalesHistoryItem {
	return SalesHistoryItem{
		ProductID: productID,
		Quantity:  quantity,
	}
}
