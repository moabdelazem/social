package main

import (
	"log"
	"net/http"
)

func (app *application) HealthChechHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "OK",
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		// TODO: More proper error handling
		log.Println("writeJSON error:", err)
	}
}
