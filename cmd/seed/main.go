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

	addr := env.GetString("DB_ADDR", "postgres://devuser:devpass@localhost:5432/myapp_dev?sslmode=disable")

	database, err := db.NewDatabase(
		addr,
		env.GetInt("DB_MAX_OPEN_CONNS", 30),
		env.GetInt("DB_MAX_IDLE_CONNS", 30),
		env.GetString("DB_MAX_IDLE_TIME", "15m"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Database connection established!")

	storage := store.NewStorage(database)

	log.Println("Starting database seeding...")
	if err := db.Seed(storage, database); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
}
