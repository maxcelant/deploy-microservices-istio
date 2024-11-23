package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxcelant/istio-microservice-sample-users/internal/svc"
)

func main() {
	lg := log.New(os.Stdout, "users-svc", log.LstdFlags)
	users, err := svc.LoadUsers()
	if err != nil {
		lg.Fatalf("Error loading JSON: %v", err)
	}

	sm := http.NewServeMux()
	sm.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	sm.Handle("/api/users", svc.UsersHandler(lg, users))
	sm.Handle("/api/user", svc.UserHandler(lg, users))

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		ErrorLog:     lg,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lg.Println("Listening on port 8080")
	lg.Fatal(s.ListenAndServe())
}
