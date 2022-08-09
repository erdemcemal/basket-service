package campaign

import (
	"github.com/erdemcemal/basket-service/internal/models"
	"github.com/shopspring/decimal"
	"math"
	"testing"
)

func TestDiscountCalculator_SameProductRule_CalculateDiscount(t *testing.T) {
	dc := NewDiscountCalculator([]Rule{SameProductRule{}})
	cart := models.ShoppingCart{
		Items: []models.ShoppingCartItem{
			{
				Quantity: 1,
				Price:    decimal.New(10, 0),
			},
			{
				Quantity: 2,
				Price:    decimal.New(10, 0),
			},
			{
				Quantity: 3,
				Price:    decimal.New(10, 0),
			},
			{
				Quantity: 4,
				Price:    decimal.New(10, 0),
			},
			{
				Quantity: 5,
				Price:    decimal.New(10, 0),
			},
		},
	}
	// there are 2 items in the cart and each item has quantity greater than 3 so discount should be applied for them.
	// 1. discount for first item is 0.8 * 10 - (item.Quantity-3) = 10 * 1 * 8 / 100 = 0.8
	// 2. discount for second item is 0.8 * 10 - (item.Quantity-3) = 10 * 2 * 8 / 100 = 1.6
	// expected result is 1.6 + 0.8 = 2.4
	expectedDiscount := 2.4
	discount := math.Round(dc.CalculateDiscount(cart)*100) / 100
	if discount != expectedDiscount {
		t.Errorf("Expected discount to be 2.4, got %f", discount)
	}
}

func TestPurchaseAmountRule_CalculateDiscount(t *testing.T) {
	cart := models.ShoppingCart{
		Items: []models.ShoppingCartItem{
			{
				Quantity: 1,
				Price:    decimal.New(10, 0),
			},
			{
				Quantity: 2,
				Price:    decimal.New(10, 0),
			},
			{
				Quantity: 3,
				Price:    decimal.New(10, 0),
			},
			{
				Quantity: 4,
				Price:    decimal.New(50, 0),
			},
			{
				Quantity: 5,
				Price:    decimal.New(50, 0),
			},
		},
	}
	cart.CalculateTotalPrice()
	purchaseAmountRule := NewPurchaseAmountRule(100, 150)
	dc := NewDiscountCalculator([]Rule{purchaseAmountRule})

	// cart total price is cart.TotalPrice = (10 * 1) + (10 * 2) + (10 * 3) + (50 * 4) + (50 * 5) = 510
	// expected discount is 510 * 10 / 100 = 51
	expectedDiscount := 51.0
	discount := math.Round(dc.CalculateDiscount(cart)*100) / 100
	if discount != expectedDiscount {
		t.Errorf("Expected discount to be 51, got %f", discount)
	}
}

func TestEveryFourthOrderRule_CalculateDiscount(t *testing.T) {
	rule := EveryFourthOrderRule{MinPurchaseAmountInMonth: 250, LastFourthOrderAmount: 350}
	dc := NewDiscountCalculator([]Rule{rule})

	cart := models.ShoppingCart{
		Items: []models.ShoppingCartItem{
			{
				Quantity: 1,
				Price:    decimal.New(10, 0),
				VatRate:  1,
			},
			{
				Quantity: 2,
				Price:    decimal.New(10, 0),
				VatRate:  1,
			},
			{
				Quantity: 3,
				Price:    decimal.New(10, 0),
				VatRate:  8,
			},
			{
				Quantity: 4,
				Price:    decimal.New(50, 0),
				VatRate:  8,
			},
			{
				Quantity: 5,
				Price:    decimal.New(50, 0),
				VatRate:  18,
			},
		},
	}
	cart.CalculateTotalPrice()

	// 10 * 3 * 0.1 = 3
	// 50 * 4 * 0.1 = 20
	// 50 * 5 * 0.15 = 37.5
	// expected discount is 37.5 + 3 + 20 = 60.5
	expectedDiscount := 60.5
	discount := math.Round(dc.CalculateDiscount(cart)*100) / 100
	if discount != expectedDiscount {
		t.Errorf("Expected discount to be 60.5, got %f", discount)
	}
}
