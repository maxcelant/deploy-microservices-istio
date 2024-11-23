package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/maxcelant/istio-microservice-sample-orders/internals/svc"
)

func main() {
	lg := log.New(os.Stdout, "orders-svc", log.LstdFlags)
	orders, err := svc.LoadOrders()
	if err != nil {
		lg.Fatalf("Error loading JSON: %v", err)
	}
	orderSvc := svc.NewOrderService(lg, orders)
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
