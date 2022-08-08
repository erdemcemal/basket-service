package http

import (
	"encoding/json"
	"github.com/erdemcemal/basket-service/internal/basket"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	Router  *mux.Router
	service basket.BasketService
	server  *http.Server
}

func NewHandler(service basket.BasketService) *Handler {
	h := &Handler{
		service: service,
	}
	h.Router = mux.NewRouter()
	h.Router.Use(JSONMiddleware)
	h.mapRoutes()

	h.server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", h.AliveCheck).Methods("GET")
	h.Router.HandleFunc("/api/v1/products", h.GetProducts).Methods("GET")
	h.Router.HandleFunc("/api/v1/basket", h.GetBasket).Methods("GET")
	h.Router.HandleFunc("/api/v1/basket", h.AddItemToBasket).Methods("POST")
	h.Router.HandleFunc("/api/v1/basket/{productId}", h.RemoveItemFromBasket).Methods("DELETE")
	h.Router.HandleFunc("/api/v1/basket/{productId}", h.UpdateItemInBasket).Methods("PUT")
	h.Router.HandleFunc("/api/v1/basket/checkout", h.CheckoutBasket).Methods("GET")
}

// AliveCheck - checks if service is alive
func (h *Handler) AliveCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "I am Alive!"}); err != nil {
		panic(err)
	}
}
