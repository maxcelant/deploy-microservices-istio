package svc

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			lg.Printf("Invalid item ID: %s", idStr)
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}

		lg.Printf("Handling request for item ID: %d", id)

		for _, item := range items {
			if item.ID == id {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(item); err != nil {
					lg.Printf("Error encoding item: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
				return
			}
		}

		lg.Printf("Item with ID %d not found", id)
		http.Error(w, "Item not found", http.StatusNotFound)
	}
}
