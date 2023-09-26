package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"log/slog"
)

type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}
func sendErr(ctx context.Context, w http.ResponseWriter, statusCode int, err error) {
	send(ctx, w, statusCode, Error{Code: "internal_error", Msg: "Internal server error"})
}

func send(ctx context.Context, w http.ResponseWriter, statusCode int, body interface{}) {
	const jsonContentType = "application/json; charset=utf-8"

	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("Unable to encode body as JSON", "error", err)
	}
}
