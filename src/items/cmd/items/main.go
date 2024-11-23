package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/maxcelant/istio-microservice-sample-items/internal/service"
)

func main() {
	lg := log.New(os.Stdout, "items ", log.LstdFlags)
	items, err := service.LoadItems()
	itemSvc := service.New(lg, items)
	if err != nil {
		lg.Fatalf("Error loading JSON: %v", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/api/items", itemSvc.GetItems()).Methods(http.MethodGet)
	router.HandleFunc("/api/items/{id}", itemSvc.GetItem()).Methods(http.MethodGet)

	s := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ErrorLog:     lg,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lg.Print("Listening on port 8081")
	lg.Fatal(s.ListenAndServe())
}
