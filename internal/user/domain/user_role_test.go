package domain_test

import (
	"testing"
	"time"

	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/stretchr/testify/assert"
)

func TestMakeModerator(t *testing.T) {
	t.Parallel()

	t.Run("should return error if user is already a moderator", func(t *testing.T) {
		t.Parallel()
		user := domain.MustNewUser("user-id", "johndoe@gmail.com",
			"johndoe", rbac.Moderator,
			time.Now(),
			time.Now(),
			nil)
		err := user.MakeModerator()
		assert.NotNil(t, err)
	})

	t.Run("should correctly change user role to moderator", func(t *testing.T) {
		t.Parallel()
		user := domain.MustNewUser(
			"user-id",
			"johndoe@gmail.com",
			"johndoe",
			rbac.Regular,
			time.Now(),
			time.Now(),
			nil)
		err := user.MakeModerator()
		assert.Nil(t, err)
		assert.Equal(t, user.Role(), rbac.Moderator)
		assert.True(t, user.IsModerator())
	})
}

func TestMakeRegular(t *testing.T) {
	t.Parallel()

	t.Run("should return error if user is already a regular", func(t *testing.T) {
		t.Parallel()
		user := domain.MustNewUser(
			"user-id",
			"johndoe@gmail.com",
			"johndoe",
			rbac.Regular,
			time.Now(),
			time.Now(),
			nil)
		err := user.MakeRegular()
		assert.NotNil(t, err)
	})

	t.Run("should correctly change user role to regular", func(t *testing.T) {
		t.Parallel()
		user := domain.MustNewUser(
			"user-id", "johndoe@gmail.com",
			"johndoe", rbac.Moderator,
			time.Now(), time.Now(),
			nil)
		err := user.MakeRegular()
		assert.Nil(t, err)
		assert.Equal(t, user.Role(), rbac.Regular)
		assert.True(t, user.IsRegular())
	})
}
