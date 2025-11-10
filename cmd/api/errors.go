package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %s, Path: %s, Method: %s", err.Error(), r.URL.Path, r.Method)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request error: %s, Path: %s, Method: %s", err.Error(), r.URL.Path, r.Method)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found error: %s, Path: %s, Method: %s", err.Error(), r.URL.Path, r.Method)
	writeJSONError(w, http.StatusNotFound, "the requested resource could not be found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Conflict error: %s, Path: %s, Method: %s", err.Error(), r.URL.Path, r.Method)
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Unauthorized error: %s, Path: %s, Method: %s", err.Error(), r.URL.Path, r.Method)
	writeJSONError(w, http.StatusUnauthorized, "you must be authenticated to access this resource")
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Unauthorized error: %s, Path: %s, Method: %s", err.Error(), r.URL.Path, r.Method)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "you must be authenticated to access this resource")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("Forbidden: Path: %s, Method: %s", r.URL.Path, r.Method)
	writeJSONError(w, http.StatusForbidden, "you do not have permission to access this resource")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	log.Printf("Rate limit exceeded: Path: %s, Method: %s", r.URL.Path, r.Method)
	w.Header().Set("Retry-After", retryAfter)
	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, please try again later")
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method not allowed: %s, Path: %s", r.Method, r.URL.Path)
	writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func (app *application) unprocessableEntityResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Unprocessable entity: %s, Path: %s, Method: %s", err.Error(), r.URL.Path, r.Method)
	writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
}
