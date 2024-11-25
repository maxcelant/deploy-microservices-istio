package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/maxcelant/istio-microservice-sample-orders/internals/cfg"
	"github.com/maxcelant/istio-microservice-sample-orders/internals/svc"
)

func initDB() (*sql.DB, func()) {
	connStr := "postgres://orders_user:orders_pass@localhost:5434/orders_db?sslmode=disable"
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

	lg := log.New(os.Stdout, "orders ", log.LstdFlags)
	config, err := cfg.LoadConfig()
	if err != nil {
		lg.Fatalf("failed to load config: %v", err)
	}
	orderService := svc.New(db, config, lg)
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")
	router.HandleFunc("/api/orders/{id}", orderService.GetOrderById()).Methods(http.MethodGet)
	router.HandleFunc("/api/orders", orderService.CreateOrder()).Methods(http.MethodPost)

	lg.Println("Starting server on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		lg.Fatalf("Could not start server: %s\n", err)
	}
}
