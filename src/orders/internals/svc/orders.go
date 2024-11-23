package svc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Status string

var (
	OpenStatus      Status = "OPEN"
	CompletedStatus Status = "COMPLETED"
	CancelledStatus Status = "CANCELLED"
)

type Order struct {
	ID         int     `json:"id"`
	UserID     int     `json:"userId"`
	Items      []int   `json:"items"`
	TotalPrice float64 `json:"totalPrice"`
	Status     Status  `json:"status"`
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
