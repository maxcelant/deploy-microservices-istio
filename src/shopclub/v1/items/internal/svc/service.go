package svc

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type Item struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ItemService struct {
	db *sql.DB
	lg *log.Logger
}

func New(db *sql.DB, lg *log.Logger) *ItemService {
	return &ItemService{
		db: db,
		lg: lg,
	}
}

func (s *ItemService) GetItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.lg.Println("Handling request for all items")

		rows, err := s.db.Query(`
			SELECT id, name, description, price, created_at 
			FROM items 
		`)
		if err != nil {
			s.lg.Printf("Error querying database: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []Item
		for rows.Next() {
			var item Item
			if err := rows.Scan(
				&item.ID, &item.Name, &item.Description, &item.Price, &item.CreatedAt,
			); err != nil {
				s.lg.Printf("Error scanning row: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}

		if err := rows.Err(); err != nil {
			s.lg.Printf("Error iterating rows: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(items); err != nil {
			s.lg.Printf("Error encoding items: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (s *ItemService) CreateItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item Item

		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			s.lg.Printf("Invalid request payload: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := s.validate(item); err != nil {
			s.lg.Printf("Validation error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO items (name, description, price, created_at)
			VALUES ($1, $2, $3, $4) RETURNING id
		`
		err := s.db.QueryRow(query,
			item.Name, item.Description, item.Price, time.Now(),
		).Scan(&item.ID)
		if err != nil {
			s.lg.Printf("Error inserting user: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(item); err != nil {
			s.lg.Printf("Error encoding user response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (s *ItemService) validate(item Item) error {
	if item.Name == "" || item.Description == "" || item.Price == 0.0 {
		return errors.New("missing required fields")
	}
	return nil
}
