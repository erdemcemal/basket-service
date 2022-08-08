package http

import (
	"encoding/json"
	"github.com/erdemcemal/basket-service/internal/dto"
	"github.com/gorilla/mux"
	"net/http"
)

// Response object for JSON responses
type Response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// GetProducts - get all products
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetProducts(r.Context())
	if err != nil {
		sendErrorResponse(w, "failed to get products", err)
		return
	}
	if err := sendOkResponse(w, products); err != nil {
		panic(err)
	}
}

// GetBasket - get basket for user with user id
func (h *Handler) GetBasket(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("user_id")
	basket, err := h.service.GetBasket(r.Context(), userId)
	if err != nil {
		sendErrorResponse(w, "failed to get basket", err)
		return
	}
	if err := sendOkResponse(w, basket); err != nil {
		panic(err)
	}
}

// AddItemToBasket - adds an item to user basket
func (h *Handler) AddItemToBasket(w http.ResponseWriter, r *http.Request) {
	var item dto.AddItemToBasketDTO
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body", err)
		return
	}
	userId := r.Header.Get("user_id")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(Response{Message: "user_id header is required"}); err != nil {
			panic(err)
		}
	}
	cart, err := h.service.AddItemToBasket(r.Context(), userId, item)
	if err != nil {
		sendErrorResponse(w, "Failed to add item to basket", err)
		return
	}
	if err := sendOkResponse(w, cart); err != nil {
		panic(err)
	}
}

// RemoveItemFromBasket - removes an item from user basket
func (h *Handler) RemoveItemFromBasket(w http.ResponseWriter, r *http.Request) {
	productId := mux.Vars(r)["productId"]
	if productId == "" {
		sendErrorResponse(w, "productId is required", nil)
		return
	}
	userId := r.Header.Get("user_id")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(Response{Message: "user_id header is required"}); err != nil {
			panic(err)
		}
	}
	cart, err := h.service.RemoveItemFromBasket(r.Context(), userId, productId)
	if err != nil {
		sendErrorResponse(w, "Failed to remove item from basket", err)
		return
	}
	if err := sendOkResponse(w, cart); err != nil {
		panic(err)
	}
}

// UpdateItemInBasket - updates an item quantity in user basket
func (h *Handler) UpdateItemInBasket(w http.ResponseWriter, r *http.Request) {
	var item dto.UpdateItemInBasketDTO
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body", err)
		return
	}
	if item.Quantity < 1 {
		sendErrorResponse(w, "Quantity must be greater than 0", nil)
		return
	}
	userId := r.Header.Get("user_id")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(Response{Message: "user_id header is required"}); err != nil {
			panic(err)
		}
	}
	cart, err := h.service.UpdateItemInBasket(r.Context(), userId, item.ProductID, item.Quantity)
	if err != nil {
		sendErrorResponse(w, "Failed to update item in basket", err)
		return
	}
	if err := sendOkResponse(w, cart); err != nil {
		panic(err)
	}
}

// CheckoutBasket - checkout user basket
func (h *Handler) CheckoutBasket(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("user_id")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(Response{Message: "user_id header is required"}); err != nil {
			panic(err)
		}
	}
	err := h.service.CheckoutBasket(r.Context(), userId)
	if err != nil {
		sendErrorResponse(w, "Failed to checkout basket", err)
		return
	}
	if err := sendOkResponse(w, nil); err != nil {
		panic(err)
	}
}

func sendOkResponse(w http.ResponseWriter, resp interface{}) error {
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
		panic(err)
	}
}
