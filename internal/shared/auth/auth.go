package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iammrsea/social-app/internal/shared/config"
)

type contextKey int

const userCtxKey contextKey = iota

type Authenticateduser struct {
	Email string
	Id    string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.Env().GoEnv() == config.Test {
			user := GetFakeUser()
			ctx := context.WithValue(r.Context(), userCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		claims, err := ParseTokenFromRequest(r)

		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		user := &Authenticateduser{
			Email: claims.Email,
		}

		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type AuthClaims struct {
	UserId string `json:"sub"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func ParseTokenFromRequest(r *http.Request) (*AuthClaims, error) {
	authHeader := r.Header.Get("Authorization")

	if strings.TrimSpace(authHeader) == "" {
		return nil, errors.New("missing Authorization header")
	}

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return nil, errors.New("invalid Authorization header")
	}

	bearerToken := parts[1]

	secret := []byte(config.Env().AuthSecret())

	token, err := jwt.ParseWithClaims(bearerToken, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*AuthClaims)

	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
