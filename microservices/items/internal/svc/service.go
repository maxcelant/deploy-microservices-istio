package svc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ItemsHandler(items []Item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func ItemHandler(items []Item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		for _, item := range items {
			if fmt.Sprintf("%d", item.ID) == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(item)
				return
			}
		}
		http.Error(w, "User not found", http.StatusNotFound)
	}
}
