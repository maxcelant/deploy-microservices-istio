package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

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

func LoadOrders() ([]Order, error) {
	file, err := os.Open("db.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file contents %v", err)
	}

	var orders []Order
	if err := json.Unmarshal(bytes, &orders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return orders, nil
}
