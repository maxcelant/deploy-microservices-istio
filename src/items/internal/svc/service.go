package svc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ItemsHandler(lg *log.Logger, items []Item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lg.Println("Handling request for all users")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(items); err != nil {
			lg.Printf("Error encoding items: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func ItemHandler(lg *log.Logger, items []Item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		lg.Printf("Handling request for item ID: %s", id)
		for _, item := range items {
			if fmt.Sprintf("%d", item.ID) == id {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(item); err != nil {
					lg.Printf("Error encoding item: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
				return
			}
		}
		lg.Printf("Item with ID %s not found", id)
		http.Error(w, "Item not found", http.StatusNotFound)
	}
}
