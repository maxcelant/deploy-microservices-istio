package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/maxcelant/istio-microservice-sample-users/internal/svc"
)

func main() {
	lg := log.New(os.Stdout, "users-svc ", log.LstdFlags)
	users, err := svc.LoadUsers()
	if err != nil {
		lg.Fatalf("Error loading JSON: %v", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/api/users", svc.UsersHandler(lg, users)).Methods(http.MethodGet)
	router.HandleFunc("/api/users/{id}", svc.UserHandler(lg, users)).Methods(http.MethodGet)

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
