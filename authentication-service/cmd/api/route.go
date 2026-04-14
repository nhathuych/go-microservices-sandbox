package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) route() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.StripSlashes)
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Post("/authenticate", app.Authenticate)

	return mux
}
