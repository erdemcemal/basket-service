package campaign

import (
	"github.com/erdemcemal/basket-service/internal/models"
	"math"
)

const (
	sameProductRuleDiscountPercentage    = 8
	purchaseAmountRuleDiscountPercentage = 10
)

// DiscountCalculator - calculates the discount for the cart with available rules
type DiscountCalculator struct {
	rules []Rule
}

// Rule - represents a rule for discount calculation
type Rule interface {
	CalculateDiscount(cart models.ShoppingCart) float64
}

// NewDiscountCalculator - creates a new discount calculator with the given rules
func NewDiscountCalculator(discountRules []Rule) *DiscountCalculator {
	return &DiscountCalculator{rules: discountRules}
}

// CalculateDiscount - calculates the discount for the given cart and return the highest discount amount
func (dc *DiscountCalculator) CalculateDiscount(cart models.ShoppingCart) float64 {
	var discount float64
	for _, rule := range dc.rules {
		discount = math.Max(rule.CalculateDiscount(cart), discount)
	}
	return discount
}
