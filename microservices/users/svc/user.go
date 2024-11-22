package svc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type User struct {
	ID          int     `json:"id"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	PhoneNumber string  `json:"phoneNumber"`
	Address     Address `json:"address"`
}

func LoadUsers() ([]User, error) {
	file, err := os.Open("db.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file contents %v", err)
	}

	var users []User
	if err := json.Unmarshal(bytes, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return users, nil
}
