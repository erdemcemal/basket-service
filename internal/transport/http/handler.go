package http

import (
	"encoding/json"
	"github.com/erdemcemal/basket-service/internal/basket"
	"github.com/gorilla/mux"
	"net/http"
)

// Handler - is a http handler for the basket service
type Handler struct {
	Router  *mux.Router
	service basket.BasketService
	server  *http.Server
}

// NewHandler - creates a new handler with the given service
func NewHandler(service basket.BasketService) *Handler {
	h := &Handler{
		service: service,
	}
	h.Router = mux.NewRouter()
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddleware)
	h.Router.Use(TimeoutMiddleware)
	h.mapRoutes()

	h.server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}
	return h
}

// Serve - starts the server
func (h *Handler) Serve() error {
	if err := h.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// mapRoutes - maps the routes to the handler
func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", h.AliveCheck).Methods("GET")
	h.Router.HandleFunc("/api/v1/products", h.GetProducts).Methods("GET")
	h.Router.HandleFunc("/api/v1/basket", Auth(h.GetBasket)).Methods("GET")
	h.Router.HandleFunc("/api/v1/basket", Auth(h.AddItemToBasket)).Methods("POST")
	h.Router.HandleFunc("/api/v1/basket/{productId}", Auth(h.RemoveItemFromBasket)).Methods("DELETE")
	h.Router.HandleFunc("/api/v1/basket", Auth(h.UpdateItemInBasket)).Methods("PUT")
	h.Router.HandleFunc("/api/v1/basket/checkout", Auth(h.CheckoutBasket)).Methods("GET")
}

// AliveCheck - checks if service is alive
func (h *Handler) AliveCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "I am Alive!"}); err != nil {
		panic(err)
	}
}
