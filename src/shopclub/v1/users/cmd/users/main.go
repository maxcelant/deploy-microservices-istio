package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/maxcelant/istio-microservice-sample-users/internal/cfg"
	"github.com/maxcelant/istio-microservice-sample-users/internal/svc"
)

func initDB(databaseURL string) (*sql.DB, func()) {
	db, err := sql.Open("postgres", databaseURL)

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
	lg := log.New(os.Stdout, "users ", log.LstdFlags)
	config, err := cfg.LoadConfig()
	if err != nil {
		lg.Fatalf("failed to load config: %v", err)
	}

	db, cleanup := initDB(config.DatabaseURL)
	defer cleanup()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
		defer cleanup()
	}

	userService := svc.New(db, lg)
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/api/users", userService.GetUsers()).Methods(http.MethodGet)
	router.HandleFunc("/api/users/{id}", userService.GetUserByID()).Methods(http.MethodGet)
	router.HandleFunc("/api/users", userService.CreateUser()).Methods(http.MethodPost)

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
