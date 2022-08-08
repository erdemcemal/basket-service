package campaign

import "github.com/erdemcemal/basket-service/internal/models"

const (
	lowVatRateDiscountPercentage  = 10
	highVatRateDiscountPercentage = 15
)

// EveryFourthOrderRule - is a rule that applies discount for every fourth order if total is more than given amount
type EveryFourthOrderRule struct {
	MinPurchaseAmountInMonth float64
	LastFourthOrderAmount    float64
}

// NewEveryFourthOrderRule - creates a new every fourth order rule with the given min purchase amount and last fourth order amount
func NewEveryFourthOrderRule(minPurchaseAmountInMonth, lastFourthOrderAmount float64) *EveryFourthOrderRule {
	return &EveryFourthOrderRule{
		MinPurchaseAmountInMonth: minPurchaseAmountInMonth,
		LastFourthOrderAmount:    lastFourthOrderAmount,
	}
}

// CalculateDiscount - calculates the discount for the given cart if user last fourth order amount totals is more than given amount
func (e EveryFourthOrderRule) CalculateDiscount(cart models.ShoppingCart) float64 {
	if e.MinPurchaseAmountInMonth >= e.LastFourthOrderAmount {
		return 0
	}
	var discount float64
	for _, item := range cart.Items {

		if item.VatRate == 8 {
			discount += item.Price.InexactFloat64() * float64(item.Quantity) * lowVatRateDiscountPercentage / 100
		}

		if item.VatRate == 18 {
			discount += item.Price.InexactFloat64() * float64(item.Quantity) * highVatRateDiscountPercentage / 100
		}
	}
	return discount
}
