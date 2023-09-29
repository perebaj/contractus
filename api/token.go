package api

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

// OAuth2Interface is an interface for the oauth2.Config struct to be able to mock it and test the callback behaviour.
type OAuth2Interface interface {
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Client(ctx context.Context, t *oauth2.Token) *http.Client
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
}

// Auth have the configuration for Auth endpoints.
type Auth struct {
	// Google OAuth2
	ClientID     string
	ClientSecret string
	RedirectURL  string

	// Default domain
	Domain string
	// JWT secret key
	JWTSecretKey string
	// Google OAuth2 config struct
	GoogleOAuthConfig OAuth2Interface
	// Access type
	AccessType string // offline(for local) or online(for production)
}

// GenerateToken generates a JWT token for the given email.
func (a *Auth) GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.JWTSecretKey))
}

// JWTAuth returns a validated JWT token.
func (a *Auth) JWTAuth() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(a.JWTSecretKey), nil)
	return tokenAuth
}
