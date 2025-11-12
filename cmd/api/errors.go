package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal server error",
		"error", err.Error(),
		"path", r.URL.Path,
		"method", r.Method,
	)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Bad request",
		"error", err.Error(),
		"path", r.URL.Path,
		"method", r.Method,
	)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Not found",
		"error", err.Error(),
		"path", r.URL.Path,
		"method", r.Method,
	)
	writeJSONError(w, http.StatusNotFound, "the requested resource could not be found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Conflict",
		"error", err.Error(),
		"path", r.URL.Path,
		"method", r.Method,
	)
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unauthorized",
		"error", err.Error(),
		"path", r.URL.Path,
		"method", r.Method,
	)
	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unauthorized (Basic Auth)",
		"error", err.Error(),
		"path", r.URL.Path,
		"method", r.Method,
	)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "you must be authenticated to access this resource")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("Forbidden",
		"path", r.URL.Path,
		"method", r.Method,
	)
	writeJSONError(w, http.StatusForbidden, "you do not have permission to access this resource")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("Rate limit exceeded",
		"path", r.URL.Path,
		"method", r.Method,
		"retry_after", retryAfter,
	)
	w.Header().Set("Retry-After", retryAfter)
	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, please try again later")
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("Method not allowed",
		"method", r.Method,
		"path", r.URL.Path,
	)
	writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func (app *application) unprocessableEntityResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unprocessable entity",
		"error", err.Error(),
		"path", r.URL.Path,
		"method", r.Method,
	)
	writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
}
