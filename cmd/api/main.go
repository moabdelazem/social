package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/moabdelazem/social/internal/db"
	"github.com/moabdelazem/social/internal/env"
	"github.com/moabdelazem/social/internal/logger"
	"github.com/moabdelazem/social/internal/store"
)

func main() {
	godotenv.Load()

	cfg := config{
		addr: env.GetString("ADDR", ":6767"),
		env:  env.GetString("ENV", "development"),
		db: dbConfig{
			addr:               env.GetString("DB_ADDR", "postgres://admin:password@localhost:5432/?sslmode=disable"),
			maxOpenConnections: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConnections: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:        env.GetString("DB_MAX_IDLE_TIME", "15m"),
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

	store := store.NewStorage(database)
	app := &application{
		config: cfg,
		store:  store,
		logger: sugar,
	}

	sugar.Infow("Application starting",
		"addr", cfg.addr,
		"env", cfg.env,
	)

	if err := app.Run(); err != nil {
		sugar.Fatalw("Failed to start server", "error", err)
	}
}
