package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

func LoadItems() ([]Item, error) {
	file, err := os.Open("db.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file contents %v", err)
	}

	var items []Item
	if err := json.Unmarshal(bytes, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return items, nil
}
