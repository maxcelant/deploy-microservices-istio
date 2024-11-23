package svc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type OrderService struct {
	lg     *log.Logger
	orders []Order
}

type OrderCreateRequest struct {
	ItemID int `json:"itemId"`
	UserID int `json:"userId"`
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
		var requestBody OrderCreateRequest
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/user?id=%d", requestBody.UserID))
		if err != nil {
			log.Fatalf("Error making GET request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Response: %s\n", string(body))

		// Using the userID, retrieve the user to verify they exist.
		// Using the itemID, retrieve the item to verify it exists
		//   And get the price of the item
		// See if the user already has an OPEN order
		//   If they do,
		//      Add this item to their list of items
		// 	 If they don't
		//      Create a new order
		//      Add the item to the list
		// Update the total order price
		//

		// o.orders = append(o.orders, newOrder)

		// file, err := json.MarshalIndent(o.orders, "", "  ")
		// if err != nil {
		// 	http.Error(w, "Error saving order", http.StatusInternalServerError)
		// 	return
		// }

		// err = os.WriteFile("db.json", file, 0644)
		// if err != nil {
		// 	http.Error(w, "Error saving order", http.StatusInternalServerError)
		// 	return
		// }

		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusCreated)
		// json.NewEncoder(w).Encode(newOrder)
	}
}
