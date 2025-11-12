package main

import (
	"errors"
	"net/http"

	"github.com/moabdelazem/social/internal/store"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	userToFollow := getUserFromCtx(r)

	// Get the authenticated user from JWT token (the follower)
	authenticatedUser := r.Context().Value("user").(*store.User)

	ctx := r.Context()
	if err := app.store.FollowerRepo.Follow(ctx, authenticatedUser.ID, userToFollow.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrorConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	app.logger.Infow("User followed",
		"follower_id", authenticatedUser.ID,
		"user_id", userToFollow.ID,
	)

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	userToUnfollow := getUserFromCtx(r)

	// Get the authenticated user from JWT token (the follower)
	authenticatedUser := r.Context().Value("user").(*store.User)

	ctx := r.Context()
	if err := app.store.FollowerRepo.Unfollow(ctx, authenticatedUser.ID, userToUnfollow.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFollowing):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	app.logger.Infow("User unfollowed",
		"follower_id", authenticatedUser.ID,
		"user_id", userToUnfollow.ID,
	)

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) getUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	ctx := r.Context()
	posts, err := app.store.PostsRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
	}
}
