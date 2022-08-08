package models

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type ShoppingCart struct {
	Base
	UserID        string             `json:"user_id"`
	Items         []ShoppingCartItem `json:"items"`
	TotalPrice    decimal.Decimal    `json:"total_price"`
	TotalVat      decimal.Decimal    `json:"total_vat"`
	TotalDiscount decimal.Decimal    `json:"total_discount"`
	SubTotal      decimal.Decimal    `json:"total_after_vat"`
}

func NewShoppingCart(userID string) ShoppingCart {
	cartId := uuid.Must(uuid.NewV4())
	return ShoppingCart{
		Base: Base{
			ID: cartId,
		},
		UserID:        userID,
		Items:         []ShoppingCartItem{},
		TotalPrice:    decimal.Zero,
		TotalDiscount: decimal.Zero,
		TotalVat:      decimal.Zero,
		SubTotal:      decimal.Zero,
	}
}

func (s *ShoppingCart) AddItem(item ShoppingCartItem) {
	s.Items = append(s.Items, item)
	s.CalculateTotalPrice()
}

func (s *ShoppingCart) RemoveItem(itemId string) {
	for i, it := range s.Items {
		if it.ProductID.String() == itemId {
			s.Items = append(s.Items[:i], s.Items[i+1:]...)
			s.CalculateTotalPrice()
			return
		}
	}
}

func (s *ShoppingCart) UpdateQuantity(item ShoppingCartItem) {
	for i, it := range s.Items {
		if it.ProductID == item.ProductID {
			s.Items[i] = item
			s.CalculateTotalPrice()
			return
		}
	}
}

func (s *ShoppingCart) ContainsItem(productId string) bool {
	for _, i := range s.Items {
		if i.ProductID.String() == productId {
			return true
		}
	}
	return false
}

func (s *ShoppingCart) GetCartItemByProductId(productId string) (ShoppingCartItem, bool) {
	for _, i := range s.Items {
		if i.ProductID.String() == productId {
			return i, true
		}
	}
	return ShoppingCartItem{}, false
}

func (s *ShoppingCart) CalculateTotalPrice() {
	// calculate total price with vat rate and quantity
	totalVat := decimal.Zero
	totalPrice := decimal.Zero
	for _, item := range s.Items {
		totalPrice = totalPrice.Add(item.Price.Mul(decimal.NewFromInt32(item.Quantity)))
		totalVat = totalVat.Add(item.Price.Mul(decimal.NewFromInt32(item.Quantity)).Mul(decimal.NewFromInt32(item.VatRate)).Div(decimal.NewFromInt32(100)))
	}
	s.TotalPrice = totalPrice
	s.TotalVat = totalVat
	s.SubTotal = decimal.Sum(totalPrice, totalVat)
}

func (s *ShoppingCart) UpdateItemQuantity(productId string, newQuantity int32) {
	for i, item := range s.Items {
		if item.ProductID.String() == productId {
			s.Items[i].Quantity = newQuantity
			s.CalculateTotalPrice()
			return
		}
	}
}
