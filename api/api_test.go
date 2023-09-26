package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO(JOJO) assert the body
func TestSendError(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	sendErr(ctx, w, http.StatusBadRequest, Error{Code: "bad_request", Msg: "Bad request"})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TODO(JOJO) assert the body
func TestSend(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	send(ctx, w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{"Hello"})
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
