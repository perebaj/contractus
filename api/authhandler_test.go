package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func TestToken(t *testing.T) {
	r := chi.NewRouter()
	RegisterAuthHandler(r, Auth{})

	wantToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9.PZLMJBT9OIVG2qgp9hQr685oVYFgRgWpcSPmNcw6y7M"
	req := httptest.NewRequest(http.MethodGet, "/token", nil)
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: wantToken,
	})

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var response struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	assert(t, resp.Code, http.StatusOK)
	assert(t, response.Token, wantToken)
}

func TestToken_Unauthorized(t *testing.T) {
	r := chi.NewRouter()
	RegisterAuthHandler(r, Auth{})

	req := httptest.NewRequest(http.MethodGet, "/token", nil)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert(t, resp.Code, http.StatusUnauthorized)
}

func TestLogin(t *testing.T) {
	a := Auth{
		GoogleOAuthConfig: &oauth2.Config{
			ClientID:     "client_id",
			ClientSecret: "client_secret",
			Endpoint:     google.Endpoint,
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		},
	}

	r := chi.NewRouter()
	RegisterAuthHandler(r, a)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert(t, resp.Code, http.StatusFound)
}
