// Package api contains the authentication endpoints.
package api

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/perebaj/contractus"
	"golang.org/x/oauth2"
)

// TODO(JOJO): randomize this
var randState = "random"

// RegisterAuthHandler register the auth endpoints.
func RegisterAuthHandler(r chi.Router, a Auth) {
	const (
		loginURL    = "/"
		callbackURL = "/callback"
		tokenURL    = "/token"
	)
	r.Method(http.MethodGet, loginURL, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login(w, r, a)
	}))
	r.Method(http.MethodGet, callbackURL, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callback(w, r, a)
	}))
	r.Method(http.MethodGet, tokenURL, http.HandlerFunc(token))
}

func login(w http.ResponseWriter, r *http.Request, a Auth) {
	var url string
	if a.AccessType == "online" {
		url = a.GoogleOAuthConfig.AuthCodeURL(randState, oauth2.AccessTypeOnline)
	} else {
		url = a.GoogleOAuthConfig.AuthCodeURL(randState, oauth2.AccessTypeOffline)
	}

	http.Redirect(w, r, url, http.StatusFound)
	slog.Info("login request received")
}

func callback(w http.ResponseWriter, r *http.Request, a Auth) {
	state := r.FormValue("state")

	if state == "" || state != randState {
		sendErr(w, http.StatusBadRequest, Error{"invalid_state", "invalid state query parameter"})
		return
	}

	code := r.FormValue("code")
	if code == "" {
		sendErr(w, http.StatusBadRequest, Error{"invalid_code", "invalid code query parameter"})
		return
	}

	ctx := r.Context()
	token, err := a.GoogleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err)
		return
	}

	client := a.GoogleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	d := json.NewDecoder(resp.Body)

	var usr contractus.GoogleUser
	err = d.Decode(&usr)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err)
		return
	}

	tokenString, err := a.GenerateToken(usr.Email)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err)
		return
	}

	cookie := http.Cookie{
		Name:   "jwt",
		Value:  tokenString,
		MaxAge: 60 * 60 * 24 * 7, // 1 week
		Domain: a.Domain,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, a.Domain+"/docs", http.StatusSeeOther)
}

func token(w http.ResponseWriter, r *http.Request) {
	jwt, err := r.Cookie("jwt")
	if err != nil {
		sendErr(w, http.StatusUnauthorized, Error{"login_again", "try to log in again"})
		return
	}

	send(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{jwt.Value})
}
