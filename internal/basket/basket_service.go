package basket

import (
	"context"
	"errors"
	"fmt"
	"github.com/erdemcemal/basket-service/internal/dto"
	"github.com/erdemcemal/basket-service/internal/models"
	basketstore "github.com/erdemcemal/basket-service/internal/store/basket"
	"gorm.io/gorm"
)

var (
	ErrGettingUserShoppingCart = errors.New("error getting user shopping cart")
	ErrProductNotFound         = errors.New("product not found")
	ErrUpdateProductQuantity   = errors.New("error updating product quantity")
	ErrProductNotInBasket      = errors.New("product not in basket")
	ErrProductStockNotEnough   = errors.New("product stock not enough")
	ErrProductAlreadyInBasket  = errors.New("product already in basket")
	ErrCheckoutBasket          = errors.New("error checking out basket")
)

type BasketService interface {
	GetProducts(ctx context.Context) ([]dto.ProductDTO, error)
	GetBasket(ctx context.Context, userId string) (dto.ShoppingCartDTO, error)
	AddItemToBasket(ctx context.Context, userId string, item dto.AddItemToBasketDTO) (dto.ShoppingCartDTO, error)
	RemoveItemFromBasket(ctx context.Context, userId string, itemToRemoveId string) (dto.ShoppingCartDTO, error)
	UpdateItemInBasket(ctx context.Context, userId string, productId string, quantity int32) (dto.ShoppingCartDTO, error)
	CheckoutBasket(ctx context.Context, userId string) error
}

type Service struct {
	store basketstore.BasketStore
}

func NewService(store basketstore.BasketStore) *Service {
	return &Service{
		store: store,
	}
}

// GetProducts - returns all products in the store
// TODO: pagination, sorting, filtering
func (s *Service) GetProducts(ctx context.Context) ([]dto.ProductDTO, error) {
	fmt.Println("getting products")
	products, err := s.store.GetProducts(ctx)
	if err != nil {
		fmt.Println("error getting products:", err)
		return []dto.ProductDTO{}, err
	}
	var dtoProducts []dto.ProductDTO
	for _, product := range products {
		dtoProducts = append(dtoProducts, fromProduct(product))
	}
	return dtoProducts, nil
}

// GetBasket - returns the shopping cart for the given user id, if not exist creates a new one
func (s *Service) GetBasket(ctx context.Context, userId string) (dto.ShoppingCartDTO, error) {
	fmt.Println("getting basket for user:", userId)
	cart, err := s.store.GetBasket(ctx, userId)
	if err != nil {
		fmt.Println("error getting basket:", err)
		return dto.ShoppingCartDTO{}, err
	}
	return fromShoppingCart(cart), nil
}

// AddItemToBasket - adds an item to the shopping cart with the given product id and quantity
func (s *Service) AddItemToBasket(ctx context.Context, userId string, item dto.AddItemToBasketDTO) (dto.ShoppingCartDTO, error) {
	shoppingCart, err := s.store.GetBasket(ctx, userId)
	if err != nil {
		return dto.ShoppingCartDTO{}, ErrGettingUserShoppingCart
	}
	if shoppingCart.ContainsItem(item.ProductID) {
		return dto.ShoppingCartDTO{}, ErrProductAlreadyInBasket
	}
	// get product from store
	product, err := s.store.GetProductById(ctx, item.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ShoppingCartDTO{}, ErrProductNotFound
		}
		return dto.ShoppingCartDTO{}, err
	}
	// check if product has enough stock
	if product.Quantity < item.Quantity {
		return dto.ShoppingCartDTO{}, ErrProductStockNotEnough
	}
	cartItem := models.NewShoppingCartItem(product.ID, product.Name, item.Quantity, product.UnitPrice, product.VatRate, shoppingCart.ID.String())
	// add item to the shopping cart
	shoppingCart.AddItem(cartItem)
	// update the shopping cart
	err = s.store.UpdateBasket(ctx, userId, shoppingCart)
	if err != nil {
		return dto.ShoppingCartDTO{}, err
	}
	return fromShoppingCart(shoppingCart), nil
}

// RemoveItemFromBasket - removes an item from the shopping cart with the given product id
func (s *Service) RemoveItemFromBasket(ctx context.Context, userId string, itemToRemoveId string) (dto.ShoppingCartDTO, error) {
	shoppingCart, err := s.store.GetBasket(ctx, userId)
	if err != nil {
		return dto.ShoppingCartDTO{}, ErrGettingUserShoppingCart
	}
	cartItemToRemove, exists := shoppingCart.GetCartItemByProductId(itemToRemoveId)
	if !exists {
		return dto.ShoppingCartDTO{}, ErrProductNotFound
	}
	shoppingCart.RemoveItem(itemToRemoveId)

	err = s.store.RemoveItemFromBasket(ctx, cartItemToRemove, shoppingCart)
	if err != nil {
		return dto.ShoppingCartDTO{}, err
	}
	return fromShoppingCart(shoppingCart), nil
}

// UpdateItemInBasket - updates the quantity of an item in the shopping cart with the given product id and new quantity
func (s *Service) UpdateItemInBasket(ctx context.Context, userId string, productId string, newQuantity int32) (dto.ShoppingCartDTO, error) {
	shoppingCart, err := s.store.GetBasket(ctx, userId)
	if err != nil {
		return dto.ShoppingCartDTO{}, ErrGettingUserShoppingCart
	}
	if !shoppingCart.ContainsItem(productId) {
		return dto.ShoppingCartDTO{}, ErrProductNotInBasket
	}
	product, err := s.store.GetProductById(ctx, productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ShoppingCartDTO{}, ErrProductNotFound
		}
		return dto.ShoppingCartDTO{}, errors.New("error getting product")
	}
	if product.Quantity < newQuantity {
		return dto.ShoppingCartDTO{}, ErrProductStockNotEnough
	}
	shoppingCart.UpdateItemQuantity(productId, newQuantity)
	err = s.store.UpdateBasket(ctx, userId, shoppingCart)
	if err != nil {
		return dto.ShoppingCartDTO{}, ErrUpdateProductQuantity
	}
	return fromShoppingCart(shoppingCart), nil
}

// CheckoutBasket - checks out the shopping cart and returns the total price of the basket
func (s *Service) CheckoutBasket(ctx context.Context, userId string) error {
	shoppingCart, err := s.store.GetBasket(ctx, userId)
	if err != nil {
		return ErrGettingUserShoppingCart
	}
	err = s.store.CheckoutBasket(ctx, shoppingCart)
	if err != nil {
		return ErrCheckoutBasket
	}
	return nil
}

// fromProduct - converts a product model to a product dto
func fromProduct(product models.Product) dto.ProductDTO {
	return dto.ProductDTO{
		ID:        product.ID.String(),
		Name:      product.Name,
		UnitPrice: product.UnitPrice,
		Quantity:  product.Quantity,
		VatRate:   product.VatRate,
	}
}

// fromShoppingCart - converts a shopping cart model to a shopping cart dto
func fromShoppingCart(cart models.ShoppingCart) dto.ShoppingCartDTO {
	return dto.ShoppingCartDTO{
		ID:            cart.ID.String(),
		UserID:        cart.UserID,
		Items:         fromShoppingCartItems(cart.Items),
		TotalPrice:    cart.TotalPrice,
		TotalVat:      cart.TotalVat,
		TotalDiscount: cart.TotalDiscount,
		TotalAfterVat: cart.SubTotal,
	}
}

// fromShoppingCartItems - converts a list of shopping cart items to a list of shopping cart items dto
func fromShoppingCartItems(items []models.ShoppingCartItem) []dto.ShoppingCartItemDTO {
	var dtoItems []dto.ShoppingCartItemDTO
	for _, item := range items {
		dtoItems = append(dtoItems, dto.ShoppingCartItemDTO{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Name:      item.ProductName,
			Price:     item.Price,
			VatRate:   item.VatRate,
			Quantity:  item.Quantity,
		})
	}
	return dtoItems
}
