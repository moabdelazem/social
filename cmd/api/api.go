package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moabdelazem/social/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	env  string
	db   dbConfig
}

type dbConfig struct {
	addr               string
	maxOpenConnections int
	maxIdleConnections int
	maxIdleTime        string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(time.Second * 67)) // SIX SEVEN

	r.Route("/v1", func(r chi.Router) {
		// Operations
		r.HandleFunc("/health", app.HealthChechHandler)

		// Posts Route Group
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Get("/", app.getPostHandler)
			})
		})
	})

	return r
}

func (a *application) Run() error {
	r := &http.Server{
		Addr:    a.config.addr,
		Handler: a.mount(),
	}
	return r.ListenAndServe()
}
