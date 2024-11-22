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
	sm.Handle("/items", svc.ItemsHandler(users))
	sm.Handle("/item", svc.ItemHandler(users))

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
