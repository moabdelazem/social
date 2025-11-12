package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/moabdelazem/social/internal/auth"
	"github.com/moabdelazem/social/internal/db"
	"github.com/moabdelazem/social/internal/env"
	"github.com/moabdelazem/social/internal/logger"
	"github.com/moabdelazem/social/internal/mailer"
	"github.com/moabdelazem/social/internal/store"
)

func main() {
	godotenv.Load()

	cfg := config{
		addr:        env.GetString("ADDR", ":6767"),
		env:         env.GetString("ENV", "development"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		db: dbConfig{
			addr:               env.GetString("DB_ADDR", "postgres://admin:password@localhost:5432/?sslmode=disable"),
			maxOpenConnections: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConnections: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:        env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		mail: mailConfig{
			smtpHost:  env.GetString("SMTP_HOST", "smtp.gmail.com"),
			smtpPort:  env.GetInt("SMTP_PORT", 587),
			smtpUser:  env.GetString("SMTP_USER", ""),
			smtpPass:  env.GetString("SMTP_PASS", ""),
			fromEmail: env.GetString("MAIL_FROM_EMAIL", "noreply@example.com"),
			exp:       time.Duration(env.GetInt("MAIL_EXPIRY_HOURS", 168)) * time.Hour, // Default 7 days (168 hours)
		},
		auth: authConfig{
			token: tokenConfig{
				secret: env.GetString("JWT_SECRET", "not-so-secret-secret"),
				exp:    time.Duration(env.GetInt("JWT_EXPIRY_HOURS", 24*7)) * time.Hour, // Default 7 days
				iss:    env.GetString("JWT_ISSUER", "social-api"),
			},
		},
	}

	// Initialize logger
	zapLogger, err := logger.New(cfg.env)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	// Create sugared logger for easier usage
	sugar := zapLogger.Sugar()

	database, err := db.NewDatabase(
		cfg.db.addr,
		cfg.db.maxOpenConnections,
		cfg.db.maxIdleConnections,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		sugar.Fatalw("Failed to connect to database",
			"error", err,
			"addr", cfg.db.addr,
		)
	}
	defer database.Close()
	sugar.Info("Database connection pool established")

	// Initialize mailer client
	smtpConfig := mailer.SMTPConfig{
		Host:     cfg.mail.smtpHost,
		Port:     cfg.mail.smtpPort,
		Username: cfg.mail.smtpUser,
		Password: cfg.mail.smtpPass,
		From:     cfg.mail.fromEmail,
	}
	mailClient := mailer.NewSMTPClient(smtpConfig)

	// Initialize JWT authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	store := store.NewStorage(database)
	app := &application{
		config:        cfg,
		store:         store,
		logger:        sugar,
		mailer:        mailClient,
		authenticator: jwtAuthenticator,
	}

	sugar.Infow("Application starting",
		"addr", cfg.addr,
		"env", cfg.env,
	)

	if err := app.Run(); err != nil {
		sugar.Fatalw("Failed to start server", "error", err)
	}
}
