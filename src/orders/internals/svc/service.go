package svc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type OrderService struct {
	lg     *log.Logger
	orders []Order
}

type OrderCreateRequest struct {
	ItemID int `json:"itemId"`
	UserID int `json:"userId"`
}

type ItemResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type UserResponse struct {
	ID          int     `json:"id"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	PhoneNumber string  `json:"phoneNumber"`
	Address     Address `json:"address"`
}

func NewOrderService(lg *log.Logger, orders []Order) *OrderService {
	return &OrderService{
		lg:     lg,
		orders: orders,
	}
}

func (s *OrderService) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		for _, order := range s.orders {
			if fmt.Sprintf("%d", order.ID) == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(order)
				return
			}
		}
		http.Error(w, "Order not found", http.StatusNotFound)
	}
}

func (s *OrderService) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody OrderCreateRequest
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user, err := s.GetUser(requestBody.UserID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving user: %v", err), http.StatusInternalServerError)
			return
		}

		item, err := s.GetItem(requestBody.ItemID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving item: %v", err), http.StatusInternalServerError)
			return
		}

		findOrder := func() (Order, bool) {
			for _, o := range s.orders {
				if o.UserID == user.ID && o.Status == "OPEN" {
					order := o
					return order, true
				}
			}
			return Order{}, false
		}

		order, ok := findOrder()
		if !ok {
			order = Order{
				ID:         len(s.orders) + 1,
				UserID:     user.ID,
				Items:      []int{item.ID},
				TotalPrice: item.Price,
				Status:     OpenStatus,
			}
		} else {
			order.Items = append(order.Items, item.ID)
			order.TotalPrice += item.Price
		}

		s.orders = append(s.orders, order)

		file, err := json.MarshalIndent(s.orders, "", "  ")
		if err != nil {
			http.Error(w, "Error saving order", http.StatusInternalServerError)
			return
		}

		err = os.WriteFile("db.json", file, 0644)
		if err != nil {
			http.Error(w, "Error saving order", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}

func (o *OrderService) GetUser(userID int) (user UserResponse, err error) {
	o.lg.Printf("Fetching user with ID: %d", userID)
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/users?id=%d", userID))
	if err != nil {
		o.lg.Printf("Error making GET request: %v", err)
		return user, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		o.lg.Printf("Error reading response body: %v", err)
		return user, fmt.Errorf("error reading response body: %v", err)
	}
	if err := json.Unmarshal(bytes, &user); err != nil {
		o.lg.Printf("Failed to unmarshal json: %v", err)
		return user, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	o.lg.Printf("Successfully fetched user: %+v", user)
	return user, nil
}

func (o *OrderService) GetItem(itemID int) (item ItemResponse, err error) {
	o.lg.Printf("Fetching item with ID: %d", itemID)
	resp, err := http.Get(fmt.Sprintf("http://localhost:8081/api/items?id=%d", itemID))
	if err != nil {
		o.lg.Printf("Error making GET request: %v", err)
		return item, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		o.lg.Printf("Error reading response body: %v", err)
		return item, fmt.Errorf("error reading response body: %v", err)
	}
	if err := json.Unmarshal(bytes, &item); err != nil {
		o.lg.Printf("Failed to unmarshal json: %v", err)
		return item, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	o.lg.Printf("Successfully fetched item: %+v", item)
	return item, nil
}