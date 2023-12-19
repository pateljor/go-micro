package main

import (
	"net/http"
	"errors"
	"fmt"
)

func (app *Config) Authenticate (w http.ResponseWriter, r *http.Request) {
	var requestPayload struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// valuidate the user agauinst the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid creds"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid creds"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}