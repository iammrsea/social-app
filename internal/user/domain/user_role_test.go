package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/stretchr/testify/assert"
)

func TestMakeModerator(t *testing.T) {
	t.Parallel()

	t.Run("should return correct error if user is already a moderator", func(t *testing.T) {
		t.Parallel()
		user := domain.MustNewUser("user-id", "johndoe@gmail.com",
			"johndoe", domain.Moderator,
			time.Now(),
			time.Now(),
			nil)
		err := user.MakeModerator()
		assert.Equal(t, err, errors.New("the user johndoe is already a moderator"))
	})

	t.Run("should correctly change user role to moderator", func(t *testing.T) {
		t.Parallel()
		user := domain.MustNewUser(
			"user-id",
			"johndoe@gmail.com",
			"johndoe",
			domain.Regular,
			time.Now(),
			time.Now(),
			nil)
		err := user.MakeModerator()
		assert.Nil(t, err)
		assert.Equal(t, user.Role(), domain.Moderator)
		assert.True(t, user.IsModerator())
	})
}

func TestMakeRegular(t *testing.T) {
	t.Parallel()

	t.Run("should return correct error if user is already a regular", func(t *testing.T) {
		t.Parallel()
		user := domain.MustNewUser(
			"user-id",
			"johndoe@gmail.com",
			"johndoe",
			domain.Regular,
			time.Now(),
			time.Now(),
			nil)
		err := user.MakeRegular()
		assert.Equal(t, err, errors.New("the user johndoe is already a regular"))
	})

	t.Run("should correctly change user role to regular", func(t *testing.T) {
		t.Parallel()
		user := domain.MustNewUser(
			"user-id", "johndoe@gmail.com",
			"johndoe", domain.Moderator,
			time.Now(), time.Now(),
			nil)
		err := user.MakeRegular()
		assert.Nil(t, err)
		assert.Equal(t, user.Role(), domain.Regular)
		assert.True(t, user.IsRegular())
	})
}
