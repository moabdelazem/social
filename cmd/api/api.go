package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	addr string
	env  string
}

func (a *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(time.Second * 67)) // SIX SEVEN

	r.HandleFunc("/health", a.HealthChechHandler)
	return r
}

func (a *application) Run() error {
	r := &http.Server{
		Addr:    a.config.addr,
		Handler: a.mount(),
	}
	return r.ListenAndServe()
}
