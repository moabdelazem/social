package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/moabdelazem/social/internal/mailer"
	"github.com/moabdelazem/social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type ActivateUserPayload struct {
	Token string `json:"token" validate:"required"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the request payload
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// Hash the password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Store the new user data with invitation
	ctx := r.Context()

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// Set invitation expiry from config
	expiry := time.Now().Add(app.config.mail.exp)

	if err := app.store.UsersRepo.CreateAndInvite(ctx, user, hashToken, expiry); err != nil {
		switch {
		case errors.Is(err, store.ErrorConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Send activation email
	activationURL := fmt.Sprintf("%s/activate?token=%s", app.config.frontendURL, plainToken)

	emailData := mailer.EmailData{
		Username:      user.Username,
		ActivationURL: activationURL,
		ExpiryTime:    expiry,
		AppName:       "Social API",
	}

	// Send email in background
	go func() {
		if _, err := app.mailer.Send(user.Email, "Activate Your Account", "user_invitation", emailData, false); err != nil {
			app.logger.Errorw("Failed to send activation email",
				"error", err,
				"email", user.Email,
			)
		} else {
			app.logger.Infow("Activation email sent",
				"email", user.Email,
				"username", user.Username,
			)
		}
	}()

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	app.logger.Infow("User registered",
		"user_id", user.ID,
		"username", user.Username,
		"email", user.Email,
	)

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload ActivateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the request payload
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.UsersRepo.Activate(ctx, payload.Token); err != nil {
		switch err {
		case store.ErrorNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	app.logger.Info("User activated successfully")

	if err := app.jsonResponse(w, http.StatusOK, map[string]string{
		"message": "User activated successfully",
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
