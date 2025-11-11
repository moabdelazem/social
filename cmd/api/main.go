package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/moabdelazem/social/internal/db"
	"github.com/moabdelazem/social/internal/env"
	"github.com/moabdelazem/social/internal/store"
)

func main() {
	godotenv.Load()

	cfg := config{
		addr: env.GetString("ADDR", ":6767"),
		env:  "development",
		db: dbConfig{
			addr:               env.GetString("DB_ADDR", "postgres://admin:password@localhost:5432/?sslmode=disable"),
			maxOpenConnections: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConnections: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:        env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	database, err := db.NewDatabase(
		cfg.db.addr,
		cfg.db.maxOpenConnections,
		cfg.db.maxIdleConnections,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatalf("Something wrong happend: %v", err)
	}
	defer database.Close()
	log.Println("database connection pool established!")

	store := store.NewStorage(database)
	app := &application{
		config: cfg,
		store:  store,
	}
	log.Printf("Application Started On %s", cfg.addr)
	log.Fatal(app.Run())
}
