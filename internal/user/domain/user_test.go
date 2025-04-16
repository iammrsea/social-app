package domain_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/stretchr/testify/assert"
)

func TestChangeUsername(t *testing.T) {
	t.Parallel()

	t.Run("should return correct error if username is empty", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		err := user.ChangeUsername("")
		assert.Equal(t, err, domain.ErrUsernameRequired)
		assert.NotEqual(t, user.Username(), "")
	})

	t.Run("should correctly change username", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		err := user.ChangeUsername("mikedoe")
		assert.Nil(t, err)
		assert.Equal(t, user.Username(), "mikedoe")
	})
}

func TestAwardBadge(t *testing.T) {
	t.Parallel()

	t.Run("should return correct error if badge is empty", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		err := user.AwardBadge("")
		badges := user.Badges()[:]
		assert.Equal(t, err, domain.ErrBadgeRequired)
		assert.Equal(t, user.Badges(), badges)
	})

	t.Run("should correctly award a badge to user", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		badge := "5 star"
		badges := user.Badges()[:]
		err := user.AwardBadge(badge)
		badges = append(badges, badge)
		assert.Nil(t, err)
		assert.Equal(t, user.Badges(), badges)
	})
}

func TestRevokeAwardedBadge(t *testing.T) {
	t.Parallel()

	t.Run("should return correct error if badge is empty", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		err := user.RevokeAwardedBadge("")
		badges := user.Badges()[:]
		assert.Equal(t, err, domain.ErrBadgeRequired)
		assert.Equal(t, user.Badges(), badges)
	})

	t.Run("should return correct error if badge was not awarded to user previously", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		badge := "5 star"
		badges := user.Badges()[:]
		err := user.RevokeAwardedBadge(badge)
		assert.Equal(t, err, fmt.Errorf("the badge %s you want to revoke hasn't been awarded to the user %s previously", badge, user.Username()))
		assert.Equal(t, user.Badges(), badges)
	})

	t.Run("should correctly revoke user's badge", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		badge := "5 star"
		err := user.AwardBadge(badge)
		assert.Nil(t, err, "user.AwardBage method is broken")
		badges := user.Badges()[:]
		err = user.RevokeAwardedBadge(badge)
		assert.Nil(t, err)
		assert.NotEqual(t, user.Badges(), badges)
		assert.False(t, slices.Contains(user.Badges(), badge))
	})
}

func TestIncrementReputationScoreBy(t *testing.T) {
	t.Parallel()

	t.Run("should return correct error if increment(0) is invalid", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		err := user.IncrementReputationScoreBy(0)
		assert.Equal(t, err, domain.ErrInvalidIncrementValue)
	})

	t.Run("should return correct error if increment(<0) is invalid", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		err := user.IncrementReputationScoreBy(-1)
		assert.Equal(t, err, domain.ErrInvalidIncrementValue)
	})

	t.Run("should correctly increase user reputation score", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		score := user.ReputationScore()
		err := user.IncrementReputationScoreBy(1)
		assert.Nil(t, err)
		assert.Equal(t, user.ReputationScore(), score+1)
	})
}

func TestDecrementReputationScoreBy(t *testing.T) {
	t.Parallel()

	t.Run("should return correct error if decrement(0) is invalid", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		err := user.DecrementReputationScoreBy(0)
		assert.Equal(t, err, domain.ErrInvalidDecrementValue)
	})

	t.Run("should return correct error if decrement(<0) is invalid", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		err := user.DecrementReputationScoreBy(-3)
		assert.Equal(t, err, domain.ErrInvalidDecrementValue)
	})

	t.Run("should correctly decrease user reputation score", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		score := user.ReputationScore()
		err := user.DecrementReputationScoreBy(1)
		assert.Nil(t, err)
		assert.Equal(t, user.ReputationScore(), score-1)
	})
}
