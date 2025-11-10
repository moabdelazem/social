package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/moabdelazem/social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	UserID  int64    `json:"user_id" validate:"required"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate required fields
	if payload.Title == "" || payload.Content == "" || payload.UserID == 0 {
		writeJSONError(w, http.StatusBadRequest, "title, content, and user_id are required fields")
		return
	}

	post := &store.Post{
		UserID:  payload.UserID,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}
	ctx := r.Context()
	if err := app.store.PostsRepo.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	post_id, err := strconv.ParseInt(idParam, 10, 64)

	ctx := r.Context()
	post, err := app.store.PostsRepo.GetByID(ctx, post_id)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
