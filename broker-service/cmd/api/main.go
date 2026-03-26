package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const serverPort = "80"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Broker service is starting on port %s\n", serverPort)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", serverPort),
		Handler:      app.route(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
