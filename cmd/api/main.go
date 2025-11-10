package main

import "log"

func main() {
	cfg := config{
		addr: ":8080",
		env:  "development",
	}
	app := &application{
		config: cfg,
	}
	log.Fatal(app.Run())
}
