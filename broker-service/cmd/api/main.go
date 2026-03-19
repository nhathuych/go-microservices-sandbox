package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const PORT = "8080"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Broker service is starting on port %s\n", PORT)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", PORT),
		Handler:      app.route(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
