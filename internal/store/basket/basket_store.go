package basket

import (
	"context"
	"errors"
	"fmt"
	"github.com/erdemcemal/basket-service/internal/models"
	"gorm.io/gorm"
)

// BasketStore - defines the interface we need our basket storage layer to implement
type BasketStore interface {
	GetProducts(ctx context.Context) ([]models.Product, error)
	GetProductById(ctx context.Context, id string) (models.Product, error)
	GetBasket(ctx context.Context, userId string) (models.ShoppingCart, error)
	UpdateBasket(ctx context.Context, userId string, newCart models.ShoppingCart) error
	RemoveItemFromBasket(ctx context.Context, cartItem models.ShoppingCartItem, newCart models.ShoppingCart) error
	CheckoutBasket(ctx context.Context, cart models.ShoppingCart) error
}

type basketStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) BasketStore {
	return &basketStore{db}
}

func (bs *basketStore) GetProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	if result := bs.db.Find(&products); result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (bs *basketStore) GetProductById(ctx context.Context, id string) (models.Product, error) {
	var product models.Product
	if result := bs.db.Where("id = ?", id).First(&product); result.Error != nil {
		return models.Product{}, result.Error
	}
	return product, nil
}

func generateShoppingCart(userId string) models.ShoppingCart {
	shoppingCart := models.NewShoppingCart(userId)
	return shoppingCart
}

func (bs *basketStore) GetBasketByUserId(ctx context.Context, userId string) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	result := bs.db.Where("user_id = ?", userId).Preload("Items").First(&cart)
	if result.Error != nil {
		return models.ShoppingCart{}, result.Error
	}
	return cart, nil
}
func (bs *basketStore) GetBasket(ctx context.Context, userId string) (models.ShoppingCart, error) {
	cart, err := bs.GetBasketByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = generateShoppingCart(userId)
			if result := bs.db.Save(&cart); result.Error != nil {
				return models.ShoppingCart{}, fmt.Errorf("error creating basket for user: %s: %w", userId, result.Error)
			}
		} else {
			return models.ShoppingCart{}, errors.New("error getting basket for user: " + userId)
		}
	}
	return cart, nil
}

func (bs *basketStore) UpdateBasket(ctx context.Context, userId string, newCart models.ShoppingCart) error {
	_, err := bs.GetBasket(ctx, userId)
	if err != nil {
		return errors.New("error getting basket for user: " + userId)
	}
	if result := bs.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&newCart); result.Error != nil {
		return errors.New("error updating basket for user: " + userId)
	}
	return nil
}

func (bs *basketStore) RemoveItemFromBasket(ctx context.Context, cartItem models.ShoppingCartItem, newCart models.ShoppingCart) error {
	tx := bs.db.Begin()
	if result := tx.Delete(&cartItem); result.Error != nil {
		tx.Rollback()
		return errors.New("error deleting item from cart: " + cartItem.ID.String())
	}
	if result := tx.Save(&newCart); result.Error != nil {
		tx.Rollback()
		return errors.New("error updating cart: " + newCart.ID.String())
	}
	tx.Commit()
	return nil
}

func (bs *basketStore) CheckoutBasket(ctx context.Context, cart models.ShoppingCart) error {
	tx := bs.db.Begin()
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
			return errors.New("error updating product: " + item.ProductID.String())
		}
	}
	orderHistory := models.NewOrderHistory(cart)
	if result := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&orderHistory); result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("error creating order history: %w", result.Error)
	}
	// delete shopping_cart_items relations when deleting shopping_cart
	if result := tx.Select("Items").Delete(&cart); result.Error != nil {
		tx.Rollback()
		return errors.New("error deleting cart: " + cart.ID.String())
	}
	tx.Commit()
	return nil
}

func (bs *basketStore) productExists(productId string) (models.Product, error) {
	var product models.Product
	if result := bs.db.Where("id = ?", productId).First(&product); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Product{}, errors.New("product does not exist: " + productId)
		}
		return models.Product{}, errors.New("error getting product: " + productId)
	}
	return product, nil
}
