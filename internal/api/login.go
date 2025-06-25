package api

import (
	"encoding/json"
	"net/http"

	"github.com/James-D-Wood/blog-api/internal/httputils"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, err := httputils.DecodeBasicAuth(r)
	if err != nil {
		app.Logger.Error("LoginHandler: failed to decode basic auth", "error", err)
		httputils.RespondWithJsonError(w, "malformatted auth header", 401)
		return
	}

	// fake login based on mock user list
	token, err := app.UserService.AuthenticateUser(user, pass)
	if err != nil {
		// return
		app.Logger.Error("LoginHandler: failed to authenticate user", "error", err)
		httputils.RespondWithJsonError(w, "user does not exist or wrong password provided", 401)
		return
	}

	resp := LoginResponse{
		Token: token,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		httputils.RespondWithJsonError(w, "internal error processing request", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
