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

	// Get authenticated user from context
	user := getUserFromCtx(r)

	ctx := r.Context()

	feed, err := app.store.PostsRepo.GetUserFeed(ctx, user.ID, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
