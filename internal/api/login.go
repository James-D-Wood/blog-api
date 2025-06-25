package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func decodeBasicAuth(r *http.Request) (username, password string, err error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", "", errors.New("no Authorization header provided")
	}

	if !strings.Contains(header, "Basic ") {
		return "", "", errors.New("basic auth not detected")
	}

	b64String := strings.Trim(header, "Basic ")
	b, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		return "", "", fmt.Errorf("problem decoding auth header: %s", err)
	}

	components := strings.Split(string(b), ":")
	if len(components) != 2 {
		return "", "", fmt.Errorf("found %d components in basic auth header - expected 2", len(components))
	}

	return components[0], components[1], nil
}

// TODO: add error response bodies
func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, err := decodeBasicAuth(r)
	if err != nil {
		app.Logger.Error(err.Error())
		w.WriteHeader(401)
		return
	}

	token, err := app.UserService.AuthenticateUser(user, pass)
	if err != nil {
		// return
		app.Logger.Error(err.Error())
		w.WriteHeader(401)
		return
	}

	resp := LoginResponse{
		Token: token,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
