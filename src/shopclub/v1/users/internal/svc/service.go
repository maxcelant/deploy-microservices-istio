package svc

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			lg.Printf("Invalid user ID: %s", idStr)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		lg.Printf("Handling request for user with ID: %d", id)

		for _, user := range users {
			if user.ID == id {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(user); err != nil {
					lg.Printf("Error encoding user: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}
		}

		lg.Printf("User with ID %d not found", id)
		http.Error(w, "User not found", http.StatusNotFound)
	}
}
