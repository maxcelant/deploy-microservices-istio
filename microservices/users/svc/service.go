package svc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func UsersHandler(users []User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func UserHandler(users []User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		for _, user := range users {
			if fmt.Sprintf("%d", user.ID) == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(user)
				return
			}
		}
		http.Error(w, "User not found", http.StatusNotFound)
	}
}
