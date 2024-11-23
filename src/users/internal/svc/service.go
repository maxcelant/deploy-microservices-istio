package svc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func UsersHandler(lg *log.Logger, users []User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lg.Println("Handling request for all users")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			lg.Printf("Error encoding users: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func UserHandler(lg *log.Logger, users []User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		lg.Printf("Handling request for user with ID: %s", id)
		for _, user := range users {
			if fmt.Sprintf("%d", user.ID) == id {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(user); err != nil {
					lg.Printf("Error encoding user: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}
		}
		lg.Printf("User with ID %s not found", id)
		http.Error(w, "User not found", http.StatusNotFound)
	}
}
