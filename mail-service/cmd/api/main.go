package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct{}

const serverPort = "80"

func main() {
	app := Config{}

	log.Println("Starting mail service on port ", serverPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: app.route(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
