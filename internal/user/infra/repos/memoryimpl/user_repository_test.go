package memoryimpl_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/iammrsea/social-app/internal/user/infra/repos/memoryimpl"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("should be able to register user", func(t *testing.T) {
		t.Parallel()
		user := domain.MustNewUser("user-id", "johndoe@gmail.com",
			"johndoe", rbac.Regular,
			time.Now(),
			time.Now(),
			nil, nil)

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.Register(ctx, user)

		assert.Nil(t, err)
		savedUser, _ := memRepo.GetUserById(ctx, user.Id())
		assert.NotNil(t, savedUser)
		userReadModel := userDomainToUserReadModel(user)
		assert.Equal(t, userReadModel.Id, savedUser.Id)
	})

	t.Run("should return correct error if user already exists", func(t *testing.T) {
		t.Parallel()

		user := domain.MustNewUser("user-id", "johndoe@gmail.com", "johndoe", rbac.Regular,
			time.Now(),
			time.Now(), nil, nil)

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.Register(ctx, user)
		assert.Nil(t, err)

		// Register same user again
		err = memRepo.Register(ctx, user)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "user with email or username already exists")
	})
}

func TestMakeModerator(t *testing.T) {
	t.Parallel()

	t.Run("should be able to make user moderator", func(t *testing.T) {
		t.Parallel()

		user := domain.MustNewUser("user-id", "johndoe@gmail.com", "johndoe", rbac.Regular,
			time.Now(),
			time.Now(), nil, nil)

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.Register(ctx, user)
		assert.Nil(t, err)

		err = memRepo.MakeModerator(ctx, user.Id(), func(user *domain.User) error {
			return user.MakeModerator()
		})
		assert.Nil(t, err)

		savedUser, err := memRepo.GetUserById(ctx, user.Id())
		assert.Nil(t, err)
		assert.Equal(t, savedUser.Role, rbac.Moderator)
	})

	t.Run("should return correct error if user does not exist", func(t *testing.T) {
		t.Parallel()

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.MakeModerator(ctx, "user-id", func(user *domain.User) error {
			return user.MakeModerator()
		})
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "user with id user-id does not exist")
	})
}

func TestAwardBadge(t *testing.T) {
	t.Parallel()

	t.Run("should be able to award badge to user", func(t *testing.T) {
		t.Parallel()

		user := domain.MustNewUser("user-id", "johndoe@gmail.com", "johndoe", rbac.Regular,
			time.Now(),
			time.Now(), nil, nil)

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.Register(ctx, user)
		assert.Nil(t, err)

		err = memRepo.AwardBadge(ctx, user.Id(), func(user *domain.User) error {
			return user.AwardBadge("4 star")
		})
		assert.Nil(t, err)

		savedUser, err := memRepo.GetUserById(ctx, user.Id())
		assert.Nil(t, err)
		assert.True(t, slices.Contains(savedUser.Reputation.Badges, "4 star"))
	})

	t.Run("should return correct error if user does not exist", func(t *testing.T) {
		t.Parallel()

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.AwardBadge(ctx, "user-id", func(user *domain.User) error {
			return user.AwardBadge("5 start")
		})
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "user with id user-id does not exist")
	})
}

func TestRevokeAwardedBadge(t *testing.T) {
	t.Parallel()

	t.Run("should be able to revoke badge from user", func(t *testing.T) {
		t.Parallel()

		user := domain.MustNewUser("user-id", "johndoe@gmail.com", "johndoe", rbac.Regular,
			time.Now(),
			time.Now(), nil, nil)

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.Register(ctx, user)
		assert.Nil(t, err)

		err = memRepo.AwardBadge(ctx, user.Id(), func(user *domain.User) error {
			return user.AwardBadge("4 star")
		})
		assert.Nil(t, err)

		err = memRepo.RevokeAwardedBadge(ctx, user.Id(), func(user *domain.User) error {
			return user.RevokeAwardedBadge("4 star")
		})
		assert.Nil(t, err)

		savedUser, err := memRepo.GetUserById(ctx, user.Id())
		assert.Nil(t, err)
		assert.False(t, slices.Contains(savedUser.Reputation.Badges, "4 star"))
	})

	t.Run("should return correct error if user does not exist", func(t *testing.T) {
		t.Parallel()

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.RevokeAwardedBadge(ctx, "user-id", func(user *domain.User) error {
			return user.RevokeAwardedBadge("4 star")
		})
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "user with id user-id does not exist")
	})

}

func TestChangeUsername(t *testing.T) {
	t.Parallel()

	t.Run("should be able to change username", func(t *testing.T) {
		t.Parallel()

		user := domain.MustNewUser("user-id", "johndoe@gmail.com", "johndoe", rbac.Regular, time.Now(),
			time.Now(), nil, nil)

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.Register(ctx, user)
		assert.Nil(t, err)

		err = memRepo.ChangeUsername(ctx, user.Id(), func(user *domain.User) error {
			return user.ChangeUsername("johndoe2")
		})
		assert.Nil(t, err)

		savedUser, err := memRepo.GetUserById(ctx, user.Id())
		assert.Nil(t, err)
		assert.Equal(t, savedUser.Username, "johndoe2")
	})

	t.Run("should return correct error if user does not exist", func(t *testing.T) {
		t.Parallel()

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.ChangeUsername(ctx, "user-id", func(user *domain.User) error {
			return user.ChangeUsername("johndoe2")
		})
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "user with id user-id does not exist")
	})
}

func TestGetUserById(t *testing.T) {
	t.Parallel()

	t.Run("should be able to get user by id", func(t *testing.T) {
		t.Parallel()

		user := domain.MustNewUser("user-id", "johndoe@gmail.com", "johndoe", rbac.Regular,
			time.Now(),
			time.Now(), nil, nil)

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		err := memRepo.Register(ctx, user)
		assert.Nil(t, err)

		savedUser, err := memRepo.GetUserById(ctx, user.Id())
		assert.Nil(t, err)
		assert.Equal(t, user.Id(), savedUser.Id)
	})

	t.Run("should return correct error if user does not exist", func(t *testing.T) {
		t.Parallel()

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		_, err := memRepo.GetUserById(ctx, "user-id")
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "user with id user-id does not exist")
	})

}

func TestGetUsers(t *testing.T) {
	t.Parallel()

	memRepo := memoryimpl.NewUserRepository(context.Background())

	ctx := context.Background()
	user := domain.MustNewUser("user-id", "johndoe@gmail.com", "johndoe", rbac.Regular,
		time.Now(),
		time.Now(), nil, nil)
	err := memRepo.Register(ctx, user)
	assert.Nil(t, err)

	users, hasNext, err := memRepo.GetUsers(ctx, domain.GetUsersOptions{})
	assert.Nil(t, err)
	assert.False(t, hasNext)
	assert.Equal(t, len(users), 1)
	assert.Equal(t, user.Id(), users[0].Id)
}

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()

	t.Run("should be able to get user by email", func(t *testing.T) {
		t.Parallel()

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()
		user := domain.MustNewUser("user-id", "johndoe@gmail.com",
			"johndoe", rbac.Regular,
			time.Now(),
			time.Now(),
			nil,
			nil)
		err := memRepo.Register(ctx, user)
		assert.Nil(t, err)

		savedUser, err := memRepo.GetUserByEmail(ctx, user.Email())
		assert.Nil(t, err)
		assert.Equal(t, user.Email(), savedUser.Email)
	})

	t.Run("should correct error if user does not exist", func(t *testing.T) {
		t.Parallel()

		memRepo := memoryimpl.NewUserRepository(context.Background())

		ctx := context.Background()

		_, err := memRepo.GetUserByEmail(ctx, "johndoe@gmail.com")
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "user with email johndoe@gmail.com does not exist")
	})
}

func userDomainToUserReadModel(user domain.User) domain.UserReadModel {
	return domain.UserReadModel{
		Username: user.Username(),
		Email:    user.Email(),
		Role:     user.Role(),
		Id:       user.Id(),
		Reputation: domain.UserReputation{
			ReputationScore: user.ReputationScore(),
			Badges:          user.Badges(),
		},
		BanStatus: domain.BanStatus{
			IsBanned:        user.IsBanned(),
			BannedAt:        user.BannedAt(),
			BanStartDate:    user.BanStartDate(),
			BanEndDate:      user.BanEndDate(),
			ReasonForBan:    user.ReasonForBan(),
			IsBanIndefinite: user.IsBanIndefinite(),
		},
		CreatedAt: user.JoinedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}
