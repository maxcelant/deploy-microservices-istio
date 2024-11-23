package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxcelant/istio-microservice-sample-items/internal/svc"
)

func main() {
	lg := log.New(os.Stdout, "items-svc", log.LstdFlags)
	users, err := svc.LoadItems()
	if err != nil {
		lg.Fatalf("Error loading JSON: %v", err)
	}

	sm := http.NewServeMux()
	sm.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	sm.Handle("/api/items", svc.ItemsHandler(lg, users))
	sm.Handle("/api/item", svc.ItemHandler(lg, users))

	s := &http.Server{
		Addr:         ":8081",
		Handler:      sm,
		ErrorLog:     lg,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lg.Print("Listening on port 8081")
	lg.Fatal(s.ListenAndServe())
}
