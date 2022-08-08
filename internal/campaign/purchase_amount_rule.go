package campaign

import "github.com/erdemcemal/basket-service/internal/models"

// PurchaseAmountRule - is a rule applies discount if the given amount is more than customer purchase amount in a month
type PurchaseAmountRule struct {
	MinPurchaseAmountInMonth      float64
	CustomerPurchaseAmountInMonth float64
}

// NewPurchaseAmountRule - creates a new purchase amount rule with the given min purchase amount and customer purchase amount
func NewPurchaseAmountRule(minPurchaseAmountInMonth, customerPurchaseAmountInMonth float64) *PurchaseAmountRule {
	return &PurchaseAmountRule{
		MinPurchaseAmountInMonth:      minPurchaseAmountInMonth,
		CustomerPurchaseAmountInMonth: customerPurchaseAmountInMonth,
	}
}

// CalculateDiscount - calculates the discount if the given amount is more than customer purchase amount in a month
func (p PurchaseAmountRule) CalculateDiscount(cart models.ShoppingCart) float64 {
	if p.MinPurchaseAmountInMonth >= p.CustomerPurchaseAmountInMonth {
		return 0
	}
	// apply purchase amount rule for total amount of the cart.
	cartAmount := cart.TotalPrice
	discountAmount := cartAmount.InexactFloat64() * purchaseAmountRuleDiscountPercentage / 100
	return discountAmount
}
