package svc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OrderService struct {
	lg     *log.Logger
	orders []Order
}

func NewOrderService(lg *log.Logger, orders []Order) *OrderService {
	return &OrderService{
		lg:     lg,
		orders: orders,
	}
}

func (o *OrderService) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		for _, order := range o.orders {
			if fmt.Sprintf("%d", order.ID) == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(order)
				return
			}
		}
		http.Error(w, "Order not found", http.StatusNotFound)
	}
}

func (o *OrderService) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}
}
