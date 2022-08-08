package campaign

import "github.com/erdemcemal/basket-service/internal/models"

// SameProductRule - represents a rule if any item quantity is more than 3 than apply the discount
type SameProductRule struct {
}

// CalculateDiscount - represents a rule if any item quantity is more than 3 than apply the discount
func (r SameProductRule) CalculateDiscount(cart models.ShoppingCart) float64 {
	var total float64
	// loop through the cart and calculate the discount for each item if item quantity is greater than 3 apply discount fourth and subsequent items.
	for _, item := range cart.Items {
		if item.Quantity > 3 {
			total += item.Price.InexactFloat64() * float64(item.Quantity-3) * sameProductRuleDiscountPercentage / 100
		}
	}
	return total
}
