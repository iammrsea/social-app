package auth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/stretchr/testify/assert"
)

func TestParseTokenFromRequest(t *testing.T) {
	t.Parallel()

	t.Run("should return zero AuthClaims if token is invalid", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/some-url", nil)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "token"))

		claims := auth.ParseTokenFromRequest(req)

		assert.NotNil(t, claims)
		assert.True(t, claims.IsZero())

	})
	t.Run("should return claims if token is valid", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/some-url", nil)
		fakeUser := auth.GetFakeUser(rbac.Moderator)
		token := auth.GenerateTestToken(fakeUser)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		claims := auth.ParseTokenFromRequest(req)

		assert.NotNil(t, claims)
		assert.False(t, claims.IsZero())
		assert.Equal(t, fakeUser.Email, claims.Email)
		assert.Equal(t, fakeUser.Role, claims.Role)
		assert.Equal(t, fakeUser.Id, claims.UserId)
		assert.False(t, claims.ExpiresAt.IsZero())
	})
}
