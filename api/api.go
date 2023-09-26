// Package api gather generic functions to deal with the API.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"log/slog"
)

// Error represents an error returned by the API.
type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}
func sendErr(w http.ResponseWriter, statusCode int, err error) {
	if statusCode >= 500 {
		slog.Error("Internal server error", "error", err)
	}
	err = Error{Code: "internal_server_error error", Msg: "Internal server error"}

	send(w, statusCode, err)
}

func send(w http.ResponseWriter, statusCode int, body interface{}) {
	const jsonContentType = "application/json; charset=utf-8"

	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("Unable to encode body as JSON", "error", err)
	}
}
