package auth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/stretchr/testify/assert"
)

func TestParseTokenFromRequest(t *testing.T) {
	t.Parallel()

	t.Run("should return error if token is invalid", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/some-url", nil)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "token"))

		claims, err := auth.ParseTokenFromRequest(req)

		assert.Nil(t, claims)
		assert.NotNil(t, err)

	})
	t.Run("should return claims if token is valid", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/some-url", nil)
		fakeUser := auth.GetFakeUser()
		token := auth.GenerateTestToken(fakeUser)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		claims, err := auth.ParseTokenFromRequest(req)

		assert.Nil(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, fakeUser.Email, claims.Email)
		assert.Equal(t, fakeUser.Id, claims.UserId)
		assert.False(t, claims.ExpiresAt.IsZero())
	})
}
