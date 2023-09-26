package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO(JOJO) assert the body
func TestSendError(t *testing.T) {
	w := httptest.NewRecorder()
	sendErr(w, http.StatusBadRequest, Error{Code: "bad_request", Msg: "Bad request"})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TODO(JOJO) assert the body
func TestSend(t *testing.T) {
	w := httptest.NewRecorder()
	send(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{"Hello"})
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
