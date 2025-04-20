package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iammrsea/social-app/internal/shared/config"
	"github.com/iammrsea/social-app/internal/shared/rbac"
)

type contextKey int

const userCtxKey contextKey = iota

type AuthenticatedUser struct {
	Email string
	Id    string
	Role  rbac.UserRole
}

func (a *AuthenticatedUser) IsZero() bool {
	return *a == AuthenticatedUser{}
}

func (a *AuthenticatedUser) IsAuthenticated() bool {
	return !a.IsZero()
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.Env().GoEnv() == config.Test {
			user := GetFakeUser(rbac.Admin)
			ctx := context.WithValue(r.Context(), userCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		claims := ParseTokenFromRequest(r)

		user := &AuthenticatedUser{}

		if !claims.IsZero() {
			user.Email = claims.Email
			user.Id = claims.UserId
			user.Role = claims.Role
		}

		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type AuthClaims struct {
	UserId string        `json:"sub"`
	Email  string        `json:"email"`
	Role   rbac.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func (c *AuthClaims) IsZero() bool {
	return c.Email == "" || c.UserId == "" || c.Role == rbac.Guest
}

func ParseTokenFromRequest(r *http.Request) *AuthClaims {
	authHeader := r.Header.Get("Authorization")
	zeroClaims := &AuthClaims{
		Role: rbac.Guest,
	}

	if strings.TrimSpace(authHeader) == "" {
		return zeroClaims
	}

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return zeroClaims
	}

	bearerToken := parts[1]

	secret := []byte(config.Env().AuthSecret())

	token, err := jwt.ParseWithClaims(bearerToken, &AuthClaims{}, func(t *jwt.Token) (any, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return zeroClaims
	}

	claims, ok := token.Claims.(*AuthClaims)

	if !ok {
		return zeroClaims
	}

	return claims
}

func GetUserFromCtx(ctx context.Context) *AuthenticatedUser {
	user, _ := ctx.Value(userCtxKey).(*AuthenticatedUser)
	return user
}
