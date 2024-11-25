package svc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/maxcelant/istio-microservice-sample-orders/internals/cfg"
)

type OrderService struct {
	db  *sql.DB
	lg  *log.Logger
	cfg *cfg.Config
}

type Status string

var (
	PendingStatus   Status = "PENDING"
	CompletedStatus Status = "COMPLETED"
	CancelledStatus Status = "CANCELLED"
)

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	OrderDate time.Time `json:"orderDate"`
	Status    Status    `json:"status"`
}

type OrderItems struct {
	OrderID  int `json:"orderId"`
	ItemID   int `json:"itemId"`
	Quantity int `json:"quantity"`
}

func New(db *sql.DB, cfg *cfg.Config, lg *log.Logger) *OrderService {
	return &OrderService{
		db:  db,
		lg:  lg,
		cfg: cfg,
	}
}

func (s *OrderService) GetOrderById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.lg.Printf("Invalid order ID: %s", idStr)
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		s.lg.Printf("Handling request for order with ID: %d", id)

		var order Order
		err = s.db.QueryRow(`
			SELECT id, user_id, order_date, status
			FROM orders WHERE id = $1
		`, id).Scan(
			&order.ID, &order.UserID, &order.OrderDate, &order.Status,
		)

		if err == sql.ErrNoRows {
			s.lg.Printf("Order with ID %d not found", id)
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		} else if err != nil {
			s.lg.Printf("Error querying database: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			s.lg.Printf("Error encoding order: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (s *OrderService) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			UserID   int `json:"userId"`
			ItemID   int `json:"itemId"`
			Quantity int `json:"quantity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			s.lg.Printf("Invalid request payload: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user, err := s.GetUser(requestBody.UserID)
		if err != nil {
			s.lg.Printf("Error validating user: %v", err)
			http.Error(w, fmt.Sprintf("Error validating user: %v", err), http.StatusInternalServerError)
			return
		}

		item, err := s.GetItem(requestBody.ItemID)
		if err != nil {
			s.lg.Printf("Error validating item: %v", err)
			http.Error(w, fmt.Sprintf("Error validating item: %v", err), http.StatusInternalServerError)
			return
		}

		tx, err := s.db.Begin()
		if err != nil {
			s.lg.Printf("Error starting transaction: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var orderID int
		orderQuery := `
			INSERT INTO orders (user_id, order_date, status)
			VALUES ($1, $2, $3) RETURNING id
		`
		err = tx.QueryRow(orderQuery, user.ID, time.Now(), "PENDING").Scan(&orderID)
		if err != nil {
			tx.Rollback()
			s.lg.Printf("Error inserting order: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		orderItemQuery := `
			INSERT INTO order_items (order_id, item_id, quantity)
			VALUES ($1, $2, $3)
		`
		_, err = tx.Exec(orderItemQuery, orderID, item.ID, requestBody.Quantity)
		if err != nil {
			tx.Rollback()
			s.lg.Printf("Error inserting order item: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			s.lg.Printf("Error committing transaction: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		order := Order{
			ID:        orderID,
			UserID:    user.ID,
			OrderDate: time.Now(),
			Status:    PendingStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(order); err != nil {
			s.lg.Printf("Error encoding order response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

type UserResponse struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

func (o *OrderService) GetUser(userID int) (user UserResponse, err error) {
	o.lg.Printf("Fetching user with ID: %d", userID)
	resp, err := http.Get(fmt.Sprintf("%s/api/users/%d", o.cfg.UserServiceURL, userID))
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

type ItemResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (o *OrderService) GetItem(itemID int) (item ItemResponse, err error) {
	o.lg.Printf("Fetching item with ID: %d", itemID)
	resp, err := http.Get(fmt.Sprintf("%s/api/items/%d", o.cfg.ItemServiceURL, itemID))
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
