package api

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

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
	GoogleOAuthConfig *oauth2.Config
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
