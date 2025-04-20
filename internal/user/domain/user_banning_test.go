package domain_test

import (
	"testing"
	"time"

	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/stretchr/testify/assert"
)

func TestBan(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		reason         string
		isInDefinitely bool
		timeline       *domain.BanTimeline
		expectedError  error
	}{
		{
			name:           "should return correct error if reason is not provided",
			reason:         "",
			isInDefinitely: false,
			timeline:       domain.NewBanTimeline(time.Now(), time.Now().Add(time.Hour*24)),
			expectedError:  domain.ErrEmptyReason},
		{
			name:           "should return correct error if ban timeline is not provided and ban is not indefinite",
			reason:         "abuse",
			isInDefinitely: false,
			timeline:       nil,
			expectedError:  domain.ErrBanTimelineRequired},
		{
			name:           "should ban a user indefinitely",
			reason:         "abuse",
			isInDefinitely: true,
			timeline:       nil,
			expectedError:  nil},
		{
			name:           "should ban a user within a timeout and not indefinitely",
			reason:         "abuse",
			isInDefinitely: false,
			timeline:       domain.NewBanTimeline(time.Now(), time.Now().Add(time.Hour*24)),
			expectedError:  nil},
		{
			name:           "should return correct error if user is already banned",
			reason:         "abuse",
			isInDefinitely: false,
			timeline:       domain.NewBanTimeline(time.Now(), time.Now().Add(time.Hour*24)),
			expectedError:  nil},
	}
	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			user := createUser()
			err := user.Ban(c.reason, c.isInDefinitely, c.timeline)

			if c.name == "should ban a user indefinitely" {
				assert.Nil(err)
				assert.True(user.IsBanned())
				assert.True(user.IsBanIndefinitely())
				assert.True(user.BanStartDate().IsZero())
				assert.True(user.BanEndDate().IsZero())
				assert.Equal(c.reason, user.ReasonForBan())
			} else if c.name == "should ban a user within a time period and not indefinitely" {
				assert.Nil(err)
				assert.True(user.IsBanned())
				assert.False(user.IsBanIndefinitely())
				assert.Equal(c.timeline, domain.NewBanTimeline(user.BanStartDate(), user.BanEndDate()))
				assert.Equal(c.reason, user.ReasonForBan())
			} else if c.name == "should return correct error if user is already banned" {
				err = user.Ban(c.reason, c.isInDefinitely, c.timeline)
				assert.Equal(err, domain.ErrUserAlreadyBanned)
				assert.True(user.IsBanned())
			} else {
				assert.Equal(err, c.expectedError, "Invalid error message")
			}
		})
	}
}

func TestUnBan(t *testing.T) {
	t.Parallel()
	t.Run("should return correct error if the user was not previously banned", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		err := user.UnBan()
		assert.Equal(t, err, domain.ErrUserIsNotBanned)
	})

	t.Run("should correctly unban a user", func(t *testing.T) {
		t.Parallel()
		user := createUser()
		mustBan(&user)
		err := user.UnBan()

		assert.Nil(t, err)
		assert.False(t, user.IsBanned())
	})

}

func mustBan(user *domain.User) {
	err := user.Ban("abuse", true, nil)
	if err != nil {
		panic("Unable to ban user")
	}
}

func createUser() domain.User {
	return domain.MustNewUser("user-id", "example@gmail.com",
		"john-doe", rbac.Regular,
		time.Now(),
		time.Now(),
		nil)
}
