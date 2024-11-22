package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxcelant/istio-microservice-sample-users/svc"
)

func main() {
	users, err := svc.LoadUsers()
	if err != nil {
		log.Fatalf("Error loading JSON: %v", err)
	}

	lg := log.New(os.Stdout, "user-svc", log.LstdFlags)
	sm := http.NewServeMux()
	sm.Handle("/users", svc.UsersHandler(users))
	sm.Handle("/user", svc.UserHandler(users))

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		ErrorLog:     lg,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Listening on port 8080")
	log.Fatal(s.ListenAndServe())
}
