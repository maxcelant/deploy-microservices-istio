package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ItemService struct {
	lg    *log.Logger
	items []Item
}

func New(lg *log.Logger, items []Item) *ItemService {
	return &ItemService{
		lg:    lg,
		items: items,
	}
}

func (s *ItemService) GetItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.lg.Println("Handling request for all users")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.items); err != nil {
			s.lg.Printf("Error encoding items: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func (s *ItemService) GetItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.lg.Printf("Invalid item ID: %s", idStr)
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}

		s.lg.Printf("Handling request for item ID: %d", id)

		for _, item := range s.items {
			if item.ID == id {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(item); err != nil {
					s.lg.Printf("Error encoding item: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
				return
			}
		}

		s.lg.Printf("Item with ID %d not found", id)
		http.Error(w, "Item not found", http.StatusNotFound)
	}
}
