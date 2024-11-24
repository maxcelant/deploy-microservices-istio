package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/maxcelant/istio-microservice-sample-users/internal/svc"
)

func initDB() (*sql.DB, func()) {
	connStr := "postgres://users_user:users_pass@localhost:5432/users_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	cleanup := func() {
		db.Close()
	}

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		defer cleanup()
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		defer cleanup()
	}

	return db, cleanup
}

func main() {

	db, cleanup := initDB()
	defer cleanup()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
		defer cleanup()
	}

	lg := log.New(os.Stdout, "users ", log.LstdFlags)
	userService := svc.New(db, lg)
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/api/users", userService.GetUsers()).Methods(http.MethodGet)
	router.HandleFunc("/api/users/{id}", userService.GetUserByID()).Methods(http.MethodGet)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     lg,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lg.Println("Listening on port 8080")
	lg.Fatal(s.ListenAndServe())
}
