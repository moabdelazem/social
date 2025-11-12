package main

import (
	"net/http"

	"github.com/moabdelazem/social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Parse pagination query parameters
	fq := store.PaginatedFeedQuery{}
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate pagination parameters
	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.PostsRepo.GetUserFeed(ctx, int64(2), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
