package svc

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserService struct {
	db *sql.DB
	lg *log.Logger
}

type User struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	CreatedAt   string `json:"createdAt"`
}

func New(db *sql.DB, lg *log.Logger) *UserService {
	return &UserService{
		db: db,
		lg: lg,
	}
}

func (s *UserService) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.lg.Println("Handling request for all users")

		rows, err := s.db.Query(`
			SELECT id, first_name, last_name, email, username, password, phone_number, address, created_at
			FROM users
		`)
		if err != nil {
			s.lg.Printf("Error querying database: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var user User
			if err := rows.Scan(
				&user.ID, &user.FirstName, &user.LastName, &user.Email,
				&user.Username, &user.Password, &user.PhoneNumber, &user.Address, &user.CreatedAt,
			); err != nil {
				s.lg.Printf("Error scanning row: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			s.lg.Printf("Error iterating rows: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			s.lg.Printf("Error encoding users: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (s *UserService) GetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.lg.Printf("Invalid user ID: %s", idStr)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		s.lg.Printf("Handling request for user with ID: %d", id)

		var user User
		err = s.db.QueryRow(`
			SELECT id, first_name, last_name, email, username, password, phone_number, address, created_at
			FROM users WHERE id = $1
		`, id).Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Email,
			&user.Username, &user.Password, &user.PhoneNumber, &user.Address, &user.CreatedAt,
		)

		if err == sql.ErrNoRows {
			s.lg.Printf("User with ID %d not found", id)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if err != nil {
			s.lg.Printf("Error querying database: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			s.lg.Printf("Error encoding user: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
