package basket

import (
	"context"
	"errors"
	"fmt"
	"github.com/erdemcemal/basket-service/internal/models"
	"gorm.io/gorm"
	"time"
)

// BasketStore - defines the interface we need our basket storage layer to implement
type BasketStore interface {
	GetProducts(ctx context.Context) ([]models.Product, error)
	GetProductById(ctx context.Context, id string) (models.Product, error)
	GetBasket(ctx context.Context, userId string) (models.ShoppingCart, error)
	UpdateBasket(ctx context.Context, userId string, newCart models.ShoppingCart) error
	RemoveItemFromBasket(ctx context.Context, cartItem models.ShoppingCartItem, newCart models.ShoppingCart) error
	CheckoutBasket(ctx context.Context, cart models.ShoppingCart) error
	GetUserMonthlyOrderAmount(ctx context.Context, userId string) (float64, error)
	GetEveryFourthOrderAmount(ctx context.Context) (float64, error)
}

type basketStore struct {
	db *gorm.DB
}

// NewBasketStore - creates a new basket store instance with the given database connection
func NewBasketStore(db *gorm.DB) BasketStore {
	return &basketStore{db}
}

// GetProducts - returns all products in the database s
func (bs *basketStore) GetProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	if result := bs.db.WithContext(ctx).Find(&products); result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

// GetProductById - returns a product with the given id
func (bs *basketStore) GetProductById(ctx context.Context, id string) (models.Product, error) {
	var product models.Product
	if result := bs.db.WithContext(ctx).Where("id = ?", id).First(&product); result.Error != nil {
		return models.Product{}, result.Error
	}
	return product, nil
}

// generateShoppingCart - generates a new shopping cart for the given user
func generateShoppingCart(userId string) models.ShoppingCart {
	shoppingCart := models.NewShoppingCart(userId)
	return shoppingCart
}

// GetBasketByUserId - returns the shopping cart for the given user
func (bs *basketStore) GetBasketByUserId(ctx context.Context, userId string) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	result := bs.db.WithContext(ctx).Where("user_id = ?", userId).Preload("Items").First(&cart)
	if result.Error != nil {
		return models.ShoppingCart{}, result.Error
	}
	return cart, nil
}

// GetBasket - returns the shopping cart for the given user, if not exists, creates a new one
func (bs *basketStore) GetBasket(ctx context.Context, userId string) (models.ShoppingCart, error) {
	cart, err := bs.GetBasketByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = generateShoppingCart(userId)
			if result := bs.db.WithContext(ctx).Save(&cart); result.Error != nil {
				return models.ShoppingCart{}, fmt.Errorf("error creating basket for user: %s: %w", userId, result.Error)
			}
		} else {
			return models.ShoppingCart{}, err
		}
	}
	return cart, nil
}

// UpdateBasket - updates the shopping cart for the given user and returns the new shopping cart
func (bs *basketStore) UpdateBasket(ctx context.Context, userId string, newCart models.ShoppingCart) error {
	_, err := bs.GetBasket(ctx, userId)
	if err != nil {
		return err
	}
	if result := bs.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Updates(&newCart); result.Error != nil {
		return err
	}
	return nil
}

// RemoveItemFromBasket - removes the given item from the shopping cart
func (bs *basketStore) RemoveItemFromBasket(ctx context.Context, cartItem models.ShoppingCartItem, newCart models.ShoppingCart) error {
	tx := bs.db.WithContext(ctx).Begin()
	if result := tx.Delete(&cartItem); result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result := tx.Save(&newCart); result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	tx.Commit()
	return nil
}

// CheckoutBasket - checks out the given shopping cart and delete the shopping cart and all its items
func (bs *basketStore) CheckoutBasket(ctx context.Context, cart models.ShoppingCart) error {
	tx := bs.db.WithContext(ctx).Begin()
	for _, item := range cart.Items {
		product, err := bs.productExists(item.ProductID.String())
		if err != nil {
			tx.Rollback()
			return err
		}
		if product.Quantity < item.Quantity {
			tx.Rollback()
			return errors.New("not enough stock for product: " + item.ProductID.String())
		}
		product.Quantity -= item.Quantity
		if result := tx.Save(&product); result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}
	orderHistory := models.NewSalesHistory(cart)
	if result := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&orderHistory); result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	// delete shopping_cart_items relations when deleting shopping_cart
	if result := tx.Select("Items").Delete(&cart); result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	tx.Commit()
	return nil
}

// productExists - checks if the product with the given id exists
func (bs *basketStore) productExists(productId string) (models.Product, error) {
	var product models.Product
	if result := bs.db.Where("id = ?", productId).First(&product); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Product{}, errors.New("product does not exist: " + productId)
		}
		return models.Product{}, result.Error
	}
	return product, nil
}

// GetUserMonthlyOrderAmount - returns the total amount of orders for the given user in a month
func (bs *basketStore) GetUserMonthlyOrderAmount(ctx context.Context, userId string) (float64, error) {
	var orders []models.SalesHistory
	if result := bs.db.WithContext(ctx).Where("user_id = ? AND created_at > ?", userId, time.Now().AddDate(0, -1, 0)).Find(&orders); result.Error != nil {
		return 0, result.Error
	}
	var total float64
	for _, order := range orders {
		total += order.SubTotal.InexactFloat64()
	}
	return total, nil
}

// GetEveryFourthOrderAmount - returns the total amount of every fourth order in a month
func (bs *basketStore) GetEveryFourthOrderAmount(ctx context.Context) (float64, error) {
	var transactionCount int64
	var total float64

	if result := bs.db.WithContext(ctx).Where("created_at > ?", time.Now().AddDate(0, -1, 0)).Order("created_at desc").Count(&transactionCount); result.Error != nil {
		return 0, result.Error
	}
	// if transactionCount divisible by 4, then get last 4 orders
	var orders []models.SalesHistory
	if transactionCount%4 == 0 {
		if result := bs.db.WithContext(ctx).Where("created_at > ?", time.Now().AddDate(0, -1, 0)).Order("created_at desc").Limit(4).Find(&orders); result.Error != nil {
			return 0, result.Error
		}
		for _, order := range orders {
			total += order.SubTotal.InexactFloat64()
		}
	}
	return total, nil
}
