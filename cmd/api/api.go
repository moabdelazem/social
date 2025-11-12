package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moabdelazem/social/internal/auth"
	"github.com/moabdelazem/social/internal/mailer"
	"github.com/moabdelazem/social/internal/store"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

type config struct {
	addr        string
	env         string
	frontendURL string
	db          dbConfig
	mail        mailConfig
	auth        authConfig
}

type authConfig struct {
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type mailConfig struct {
	smtpHost  string
	smtpPort  int
	smtpUser  string
	smtpPass  string
	fromEmail string
	exp       time.Duration
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
	r.Use(app.zapLogger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(time.Second * 67)) // SIX SEVEN

	// Custom 404 Not Found handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writeJSONError(w, http.StatusNotFound, "the requested resource could not be found")
	})

	// Custom 405 Method Not Allowed handler
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowedResponse(w, r)
	})

	r.Route("/v1", func(r chi.Router) {
		// Operations
		r.HandleFunc("/health", app.HealthChechHandler)

		// Posts Route Group
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})

		// Users Route Group
		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Use(app.usersContextMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerUserHandler)
			r.Put("/activate", app.activateUserHandler)
			r.Post("/login", app.loginUserHandler)
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
