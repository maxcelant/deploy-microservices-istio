package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/maxcelant/istio-microservice-sample-orders/internals/cfg"
	"github.com/maxcelant/istio-microservice-sample-orders/internals/service"
)

func main() {
	lg := log.New(os.Stdout, "orders ", log.LstdFlags)
	config, err := cfg.LoadConfig()
	if err != nil {
		lg.Fatalf("failed to load config: %v", err)
	}
	orders, err := service.LoadOrders()
	if err != nil {
		lg.Fatalf("Error loading JSON: %v", err)
	}
	orderSvc := service.NewOrderService(lg, config, orders)
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")
	router.HandleFunc("/api/order", orderSvc.Get()).Methods("GET")
	router.HandleFunc("/api/order", orderSvc.Create()).Methods("POST")

	lg.Println("Starting server on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		lg.Fatalf("Could not start server: %s\n", err)
	}
}
