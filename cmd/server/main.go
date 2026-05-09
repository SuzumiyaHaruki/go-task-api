package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func main() {
	port := getenv("APP_PORT", "8080")

	a := newApp()
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      logRequests(a.mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("go-task-api listening on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
