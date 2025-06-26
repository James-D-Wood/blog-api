package httputils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponseBody struct {
	Error string `json:"error"`
}

func RespondWithJsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	respBody := ErrorResponseBody{
		Error: message,
	}
	respBytes, _ := json.Marshal(respBody)
	w.Write(respBytes)
}

func RespondWithJson(w http.ResponseWriter, body any) {
	respBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(respBytes)
}
