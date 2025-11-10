package main

import (
	"log"

	"github.com/moabdelazem/social/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":6767"),
		env:  "development",
	}
	app := &application{
		config: cfg,
	}
	log.Fatal(app.Run())
}
