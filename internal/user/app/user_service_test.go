package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/iammrsea/social-app/internal/shared/auth"
	guard_mocks "github.com/iammrsea/social-app/internal/shared/guards/mocks"
	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/iammrsea/social-app/internal/shared/pagination"
	service "github.com/iammrsea/social-app/internal/user/app"
	"github.com/iammrsea/social-app/internal/user/app/command"
	"github.com/iammrsea/social-app/internal/user/app/query"
	"github.com/iammrsea/social-app/internal/user/domain"
	domain_mocks "github.com/iammrsea/social-app/internal/user/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type commandTestCase[T any] struct {
	name        string
	command     T
	expectedErr error
	authUser    *auth.AuthenticatedUser
	setupMocks  func(t *testing.T, repo *domain_mocks.MockUserRepository, guards *guard_mocks.MockGuards, command *T, authUser *auth.AuthenticatedUser)
}
type queryResult[D any] struct {
	data *D
	err  error
}
type queryTestCase[T any, D any] struct {
	name           string
	query          T
	expectedResult queryResult[D]
	authUser       *auth.AuthenticatedUser
	setupMocks     func(t *testing.T, repo *domain_mocks.MockUserReadModelRepository, guards *guard_mocks.MockGuards, query T, authUser *auth.AuthenticatedUser)
}

func TestCommandHandler(t *testing.T) {
	t.Parallel()
	t.Run("Register", func(t *testing.T) {
		t.Parallel()
		testRegister(t)
	})
	t.Run("BanUser", func(t *testing.T) {
		t.Parallel()
		testBanUser(t)
	})
	t.Run("UnbanUser", func(t *testing.T) {
		t.Parallel()
		testUnbanUser(t)
	})
	t.Run("ChangeUsername", func(t *testing.T) {
		t.Parallel()
		testChangeUsername(t)
	})
	t.Run("AwardBadge", func(t *testing.T) {
		t.Parallel()
		testAwardBadge(t)
	})
	t.Run("RevokeBadge", func(t *testing.T) {
		t.Parallel()
		testRevokeBadge(t)
	})
	t.Run("MakeModerator", func(t *testing.T) {
		t.Parallel()
		testMakeModerator(t)
	})
}
func TestQueryHandler(t *testing.T) {
	t.Run("GetUserByEmail", func(t *testing.T) {
		t.Parallel()
		testGetUserByEmail(t)
	})
	t.Run("GetUserById", func(t *testing.T) {
		t.Parallel()
		testGetUserById(t)
	})
	t.Run("GetUsers", func(t *testing.T) {
		t.Parallel()
		testGetUsers(t)
	})
}

func testGetUsers(t *testing.T) {
	testCases := []queryTestCase[query.GetUsers, query.Result]{
		{
			name: "authorized user can get list of users",
			query: query.GetUsers{
				First: 10,
				After: "",
			},
			authUser: &auth.AuthenticatedUser{
				Id:    "userId-1",
				Role:  rbac.Admin,
				Email: "admin@example.com",
			},
			expectedResult: queryResult[query.Result]{
				data: &query.Result{
					Data: []*domain.UserReadModel{
						{
							Id:       "userId-1",
							Email:    "user1@example.com",
							Username: "user1",
							Role:     rbac.Regular,
						},
						{
							Id:       "userId-2",
							Email:    "user2@example.com",
							Username: "user2",
							Role:     rbac.Regular,
						},
					},
					PaginationInfo: &pagination.PagenationInfo{
						HasNext: false,
					},
				},
				err: nil,
			},
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserReadModelRepository, guards *guard_mocks.MockGuards, query query.GetUsers, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.ListUsers).Return(nil)
				repo.EXPECT().GetUsers(mock.Anything, query).RunAndReturn(
					func(ctx context.Context, opts domain.GetUsersOptions) ([]*domain.UserReadModel, bool, error) {
						result := []*domain.UserReadModel{
							{
								Id:       "userId-1",
								Email:    "user1@example.com",
								Username: "user1",
								Role:     rbac.Regular,
							},
							{
								Id:       "userId-2",
								Email:    "user2@example.com",
								Username: "user2",
								Role:     rbac.Regular,
							},
						}
						return result, false, nil
					})
			},
		},
		{
			name: "unauthorized user cannot get list of users",
			query: query.GetUsers{
				First: 10,
				After: "",
			},
			authUser: &auth.AuthenticatedUser{
				Id:    "",
				Role:  rbac.Guest,
				Email: "",
			},
			expectedResult: queryResult[query.Result]{
				data: nil,
				err:  rbac.ErrUnauthorized,
			},
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserReadModelRepository, guards *guard_mocks.MockGuards, query query.GetUsers, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.ListUsers).Return(rbac.ErrUnauthorized)
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, service := setupQueryUserService(t, tt)
			res, err := service.GetUsers.Handle(ctx, tt.query)
			assertError(t, err, tt.expectedResult.err)
			assert.Equal(t, tt.expectedResult.data, res)
		})
	}
}

func testGetUserById(t *testing.T) {
	testCases := []queryTestCase[query.GetUserById, domain.UserReadModel]{
		{
			name: "authorized user get user by id",
			authUser: &auth.AuthenticatedUser{
				Id:    "userId-1",
				Role:  rbac.Regular,
				Email: "user@example.com",
			},
			query: query.GetUserById{
				Id: "userId-1234",
			},
			expectedResult: queryResult[domain.UserReadModel]{
				data: &domain.UserReadModel{
					Id:       "userId-1234",
					Email:    "user2@example.com",
					Username: "username",
					Role:     rbac.Regular,
				},
				err: nil,
			},
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserReadModelRepository, guards *guard_mocks.MockGuards, query query.GetUserById, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.ViewUser).Return(nil)
				repo.EXPECT().GetUserById(mock.Anything, query.Id).RunAndReturn(
					func(ctx context.Context, id string) (*domain.UserReadModel, error) {
						return &domain.UserReadModel{
							Id:       "userId-1234",
							Email:    "user2@example.com",
							Username: "username",
							Role:     rbac.Regular,
						}, nil
					})
			},
		},
		{
			name: "unauthorized user cannot get user by id",
			authUser: &auth.AuthenticatedUser{
				Id:    "",
				Role:  rbac.Guest,
				Email: "",
			},
			query: query.GetUserById{
				Id: "userId-1234",
			},
			expectedResult: queryResult[domain.UserReadModel]{
				data: nil,
				err:  rbac.ErrUnauthorized,
			},
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserReadModelRepository, guards *guard_mocks.MockGuards, query query.GetUserById, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.ViewUser).Return(rbac.ErrUnauthorized)
			},
		},
		{
			name: "user not found",
			authUser: &auth.AuthenticatedUser{
				Id:    "userId-1234",
				Role:  rbac.Regular,
				Email: "user@example.com",
			},
			query: query.GetUserById{
				Id: "userId-1234",
			},
			expectedResult: queryResult[domain.UserReadModel]{
				data: nil,
				err:  domain.ErrUserNotFound,
			},
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserReadModelRepository, guards *guard_mocks.MockGuards, query query.GetUserById, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.ViewUser).Return(nil)
				repo.EXPECT().GetUserById(mock.Anything, query.Id).RunAndReturn(
					func(ctx context.Context, id string) (*domain.UserReadModel, error) {
						return nil, domain.ErrUserNotFound
					})
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, service := setupQueryUserService(t, tt)
			res, err := service.GetUserById.Handle(ctx, tt.query)
			assertError(t, err, tt.expectedResult.err)
			assert.Equal(t, tt.expectedResult.data, res)
		})
	}
}

func testGetUserByEmail(t *testing.T) {
	testCases := []queryTestCase[query.GetUserByEmail, domain.UserReadModel]{
		{
			name: "authorized user get user by email",
			authUser: &auth.AuthenticatedUser{
				Id:    "userId-1234",
				Role:  rbac.Regular,
				Email: "user@example.com",
			},
			query: query.GetUserByEmail{
				Email: "user2@example.com",
			},
			expectedResult: queryResult[domain.UserReadModel]{
				data: &domain.UserReadModel{
					Id:       "userId-1234",
					Email:    "user2@example.com",
					Username: "username",
					Role:     rbac.Regular,
				},
				err: nil,
			},
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserReadModelRepository, guards *guard_mocks.MockGuards, query query.GetUserByEmail, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.ViewUser).Return(nil)
				repo.EXPECT().GetUserByEmail(mock.Anything, query.Email).RunAndReturn(
					func(ctx context.Context, email string) (*domain.UserReadModel, error) {
						return &domain.UserReadModel{
							Id:       "userId-1234",
							Email:    "user2@example.com",
							Username: "username",
							Role:     rbac.Regular,
						}, nil
					})
			},
		},
		{
			name: "unauthorized user cannot get user by email",
			authUser: &auth.AuthenticatedUser{
				Id:    "",
				Role:  rbac.Guest,
				Email: "",
			},
			query: query.GetUserByEmail{
				Email: "user2@example.com",
			},
			expectedResult: queryResult[domain.UserReadModel]{
				data: nil,
				err:  rbac.ErrUnauthorized,
			},
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserReadModelRepository, guards *guard_mocks.MockGuards, query query.GetUserByEmail, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.ViewUser).Return(rbac.ErrUnauthorized)
			},
		},
		{
			name: "user not found",
			authUser: &auth.AuthenticatedUser{
				Id:    "user-1",
				Role:  rbac.Regular,
				Email: "user1@example.com",
			},
			query: query.GetUserByEmail{
				Email: "user2@example.com",
			},
			expectedResult: queryResult[domain.UserReadModel]{
				data: nil,
				err:  domain.ErrUserNotFound,
			},
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserReadModelRepository, guards *guard_mocks.MockGuards, query query.GetUserByEmail, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.ViewUser).Return(nil)
				repo.EXPECT().GetUserByEmail(mock.Anything, query.Email).Return(nil, domain.ErrUserNotFound)
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, service := setupQueryUserService(t, tt)
			res, err := service.GetUserByEmail.Handle(ctx, tt.query)
			assertError(t, err, tt.expectedResult.err)
			assert.Equal(t, tt.expectedResult.data, res)
		})
	}
}

func testUnbanUser(t *testing.T) {
	testCases := []commandTestCase[command.UnbanUser]{
		{
			name: "authorized user can unban user",
			authUser: &auth.AuthenticatedUser{
				Id:    "userId-1234",
				Role:  rbac.Admin,
				Email: "admin@example.com",
			},
			command: command.UnbanUser{
				Id: "userId-1",
			},
			expectedErr: nil,
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserRepository, guards *guard_mocks.MockGuards, command *command.UnbanUser, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.UnbanUser).Return(nil)
				repo.EXPECT().UnbanUser(mock.Anything, command.Id, mock.AnythingOfType("func(*domain.User) error")).RunAndReturn(
					func(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
						user, err := domain.NewUser("userId-1", "user@example.com", "username", rbac.Regular,
							time.Now(), time.Now(), nil, domain.NewBan(true, "bullying", false, time.Now(), time.Now().Add(time.Hour*48), time.Now()))
						require.NoError(t, err)
						require.True(t, user.IsBanned())
						err = updateFn(&user)
						require.NoError(t, err)
						require.False(t, user.IsBanned(), "User wasn't unbanned as expected")
						return nil
					})
			},
		},
		{
			name: "unauthorized user cannot unban user",
			authUser: &auth.AuthenticatedUser{
				Id:    "userId-1234",
				Role:  rbac.Guest,
				Email: "guest@example.com",
			},
			command: command.UnbanUser{
				Id: "userId-1",
			},
			expectedErr: rbac.ErrUnauthorized,
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserRepository, guards *guard_mocks.MockGuards, command *command.UnbanUser, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.UnbanUser).Return(rbac.ErrUnauthorized)
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, userService := setupCommandUserService(t, tt)
			err := userService.UnbanUser.Handle(ctx, tt.command)
			assertError(t, err, tt.expectedErr)
		})
	}
}

func testMakeModerator(t *testing.T) {
	testCases := []commandTestCase[command.MakeModerator]{
		{
			name: "authorized user can make user moderator",
			authUser: &auth.AuthenticatedUser{
				Id:    "userId-1234",
				Role:  rbac.Admin,
				Email: "admin@example.com",
			},
			command: command.MakeModerator{
				Id: "userId-3030303048933",
			},
			expectedErr: nil,
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserRepository, guards *guard_mocks.MockGuards, command *command.MakeModerator, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.MakeModerator).Return(nil)
				repo.EXPECT().MakeModerator(mock.Anything, command.Id, mock.AnythingOfType("func(*domain.User) error")).RunAndReturn(
					func(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
						user, err := domain.NewUser("userId-239", "user@example.com", "username", rbac.Regular, time.Now(), time.Now(), nil, nil)
						require.NoError(t, err, "Unable to create new user")
						err = updateFn(&user)
						require.NoError(t, err, "Unable to make user moderator")
						require.Equal(t, user.Role(), rbac.Moderator, "User role was not updated to moderator as expected")
						return nil
					})
			},
		},
		{
			name: "unauthorized user cannot make user moderator",
			authUser: &auth.AuthenticatedUser{
				Role:  rbac.Regular,
				Id:    "userId-12350500",
				Email: "user@example.com",
			},
			command: command.MakeModerator{
				Id: "userId-3030303048933",
			},
			expectedErr: rbac.ErrUnauthorized,
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserRepository, guards *guard_mocks.MockGuards, command *command.MakeModerator, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.MakeModerator).Return(rbac.ErrUnauthorized)
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, userService := setupCommandUserService(t, tt)
			err := userService.MakeModerator.Handle(ctx, tt.command)
			assertError(t, err, tt.expectedErr)
		})
	}
}

func testRevokeBadge(t *testing.T) {
	testCases := []commandTestCase[command.RevokeAwardedBadge]{
		{
			name: "authorized user can revoke badge",
			authUser: &auth.AuthenticatedUser{
				Role:  rbac.Admin,
				Id:    "userId-12350500",
				Email: "admin@example.com",
			},
			command: command.RevokeAwardedBadge{
				Id:    "userId-123",
				Badge: "5-stars",
			},
			expectedErr: nil,
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserRepository, guards *guard_mocks.MockGuards, command *command.RevokeAwardedBadge, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.RevokeBadge).Return(nil)
				repo.EXPECT().RevokeAwardedBadge(mock.Anything, command.Id, mock.AnythingOfType("func(*domain.User) error")).RunAndReturn(
					func(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
						user, err := domain.NewUser("userId-123", "user@example.com", "username", rbac.Regular,
							time.Now(), time.Now(), domain.MustNewUserReputation(5, []string{"5-stars"}), nil)
						require.NoError(t, err, "Unable to create user")
						require.Contains(t, user.Badges(), "5-stars")
						err = updateFn(&user)
						require.NoError(t, err, "Unable to revoke badge")
						require.NotContains(t, user.Badges(), "5-stars")
						return nil
					})
			},
		},
		{
			name: "unauthorized user cannot revoke badge",
			authUser: &auth.AuthenticatedUser{
				Role:  rbac.Regular,
				Id:    "userId-123",
				Email: "user@example.com",
			},
			command: command.RevokeAwardedBadge{
				Id:    "userId-12345550505",
				Badge: "3-stars",
			},
			expectedErr: rbac.ErrUnauthorized,
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserRepository, guards *guard_mocks.MockGuards, command *command.RevokeAwardedBadge, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.RevokeBadge).Return(rbac.ErrUnauthorized)
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, userService := setupCommandUserService(t, tt)
			err := userService.RevokeAwardedBadge.Handle(ctx, tt.command)
			assertError(t, err, tt.expectedErr)
		})
	}
}

func testAwardBadge(t *testing.T) {
	testCases := []commandTestCase[command.AwardBadge]{
		{
			name: "authorized user can award badge",
			authUser: &auth.AuthenticatedUser{
				Role:  rbac.Admin,
				Id:    "userId-12350500",
				Email: "admin@example.com",
			},
			command: command.AwardBadge{
				Id:    "userId-123",
				Badge: "5-stars",
			},
			expectedErr: nil,
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserRepository, guards *guard_mocks.MockGuards, command *command.AwardBadge, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.AwardBadge).Return(nil)
				repo.EXPECT().AwardBadge(mock.Anything, command.Id, mock.AnythingOfType("func(*domain.User) error")).RunAndReturn(
					func(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
						user, err := domain.NewUser("userId-123", "user@example.com", "username", rbac.Regular, time.Now(), time.Now(), nil, nil)
						require.NoError(t, err)
						err = updateFn(&user)
						require.NoError(t, err)
						require.Contains(t, user.Badges(), command.Badge)
						return nil
					})
			},
		},
		{
			name: "unauthorized user cannot award badge",
			authUser: &auth.AuthenticatedUser{
				Role:  rbac.Regular,
				Id:    "userId-123",
				Email: "user@example.com",
			},
			command: command.AwardBadge{
				Id:    "userId-12333020202",
				Badge: "4-stars",
			},
			expectedErr: rbac.ErrUnauthorized,
			setupMocks: func(t *testing.T, repo *domain_mocks.MockUserRepository, guards *guard_mocks.MockGuards, command *command.AwardBadge, authUser *auth.AuthenticatedUser) {
				guards.EXPECT().Authorize(authUser.Role, rbac.AwardBadge).Return(rbac.ErrUnauthorized)
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, userService := setupCommandUserService(t, tt)
			err := userService.AwardBadge.Handle(ctx, tt.command)
			assertError(t, err, tt.expectedErr)
		})
	}
}

func testChangeUsername(t *testing.T) {
	testCases := []commandTestCase[command.ChangeUsername]{
		{
			name:        "user can change their username",
			expectedErr: nil,
			authUser: &auth.AuthenticatedUser{
				Role:  rbac.Regular,
				Id:    "userId-123",
				Email: "user@example.com",
			},
			command: command.ChangeUsername{
				Id:       "userId-123",
				Username: "newUsername",
			},
			setupMocks: func(t *testing.T, userRepo *domain_mocks.MockUserRepository, guard *guard_mocks.MockGuards, cmd *command.ChangeUsername, authUser *auth.AuthenticatedUser) {
				userRepo.EXPECT().UserExists(mock.Anything, "", cmd.Username).Return(false, nil)
				guard.EXPECT().CanChangeUsername(cmd.Id, authUser).Return(nil)
				userRepo.EXPECT().ChangeUsername(mock.Anything, cmd.Id, mock.AnythingOfType("func(*domain.User) error")).RunAndReturn(
					func(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
						user, err := domain.NewUser("userId-123", "user@example.com", "username", rbac.Regular, time.Now(), time.Now(), nil, nil)
						require.NoError(t, err)
						err = updateFn(&user)
						require.NoError(t, err)
						require.Equal(t, "newUsername", user.Username(), "Username was not updated as expected")
						return nil
					})

			},
		},
		{
			name:        "user cannot change their username with existing username",
			expectedErr: domain.ErrEmailOrUsernameAlreadyExists,
			authUser: &auth.AuthenticatedUser{
				Role:  rbac.Regular,
				Id:    "userId-123",
				Email: "user@example.com",
			},
			command: command.ChangeUsername{
				Id:       "userId-123",
				Username: "username",
			},
			setupMocks: func(t *testing.T, userRepo *domain_mocks.MockUserRepository, guard *guard_mocks.MockGuards, cmd *command.ChangeUsername, authUser *auth.AuthenticatedUser) {
				userRepo.EXPECT().UserExists(mock.Anything, "", cmd.Username).Return(true, nil)
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, userService := setupCommandUserService(t, tt)
			err := userService.ChangeUsername.Handle(ctx, tt.command)
			assertError(t, err, tt.expectedErr)
		})
	}
}

func testBanUser(t *testing.T) {
	testCases := []commandTestCase[command.BanUser]{
		{
			name:        "authorized user bans target user",
			expectedErr: nil,
			authUser: &auth.AuthenticatedUser{
				Id:    "userId-12345",
				Email: "admin@example.com",
				Role:  rbac.Admin,
			},
			command: command.BanUser{
				Id:             "userId-123",
				Reason:         "abuse",
				IsIndefinitely: true,
			},
			setupMocks: func(t *testing.T, userRepo *domain_mocks.MockUserRepository, guard *guard_mocks.MockGuards, cmd *command.BanUser, authUser *auth.AuthenticatedUser) {
				guard.EXPECT().Authorize(authUser.Role, rbac.BanUser).Return(nil)
				userRepo.EXPECT().BanUser(mock.Anything, cmd.Id, mock.AnythingOfType("func(*domain.User) error")).RunAndReturn(
					func(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
						user, err := domain.NewUser(cmd.Id, "testuser@gmail.com", "testuser", rbac.Regular, time.Now(), time.Now(), nil, nil)
						require.NoError(t, err, "Failed to create user")
						err = updateFn(&user)
						if err != nil {
							return err
						}
						require.True(t, user.IsBanned(), "User was not banned as expected")
						return nil
					})
			},
		},
		{
			name:        "unauthorized user cannot ban target user",
			expectedErr: rbac.ErrUnauthorized,
			authUser: &auth.AuthenticatedUser{
				Id:    "userId-12345",
				Email: "user@example.com",
				Role:  rbac.Regular,
			},
			command: command.BanUser{
				Id: "userId-123",
			},
			setupMocks: func(t *testing.T, userRepo *domain_mocks.MockUserRepository, guard *guard_mocks.MockGuards, cmd *command.BanUser, authUser *auth.AuthenticatedUser) {
				guard.EXPECT().Authorize(authUser.Role, rbac.BanUser).Return(rbac.ErrUnauthorized)
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, userService := setupCommandUserService(t, tt)
			err := userService.BanUser.Handle(ctx, tt.command)
			assertError(t, err, tt.expectedErr)
		})
	}
}

func testRegister(t *testing.T) {
	testCases := []commandTestCase[command.RegisterUser]{
		{
			name:        "register new account",
			expectedErr: nil,
			authUser: &auth.AuthenticatedUser{
				Id:    "",
				Email: "",
				Role:  rbac.Guest,
			},
			command: command.RegisterUser{
				Email:    "test@example.com",
				Username: "testuser",
			},
			setupMocks: func(t *testing.T, userRepo *domain_mocks.MockUserRepository, guard *guard_mocks.MockGuards, cmd *command.RegisterUser, authUser *auth.AuthenticatedUser) {
				guard.EXPECT().Authorize(authUser.Role, rbac.CreateAccount).Return(nil)
				userRepo.EXPECT().UserExists(mock.Anything, cmd.Email, cmd.Username).Return(false, nil)
				userRepo.EXPECT().Register(mock.Anything, mock.AnythingOfType("domain.User")).Return(nil)
			},
		},
		{
			name:        "cannot create new account with existing email",
			expectedErr: domain.ErrEmailOrUsernameAlreadyExists,
			authUser: &auth.AuthenticatedUser{
				Id:    "",
				Email: "",
				Role:  rbac.Guest,
			},
			command: command.RegisterUser{
				Email:    "test@example.com",
				Username: "testuser",
			},
			setupMocks: func(t *testing.T, userRepo *domain_mocks.MockUserRepository, guard *guard_mocks.MockGuards, cmd *command.RegisterUser, authUser *auth.AuthenticatedUser) {
				guard.EXPECT().Authorize(authUser.Role, rbac.CreateAccount).Return(nil)
				userRepo.EXPECT().UserExists(mock.Anything, cmd.Email, cmd.Username).Return(true, nil)
			},
		},
		{
			name:        "cannot create new account with existing username",
			expectedErr: domain.ErrEmailOrUsernameAlreadyExists,
			authUser: &auth.AuthenticatedUser{
				Id:    "",
				Email: "",
				Role:  rbac.Guest,
			},
			command: command.RegisterUser{
				Email:    "test@example.com",
				Username: "testuser",
			},
			setupMocks: func(t *testing.T, userRepo *domain_mocks.MockUserRepository, guard *guard_mocks.MockGuards, cmd *command.RegisterUser, authUser *auth.AuthenticatedUser) {
				guard.EXPECT().Authorize(authUser.Role, rbac.CreateAccount).Return(nil)
				userRepo.EXPECT().UserExists(mock.Anything, cmd.Email, cmd.Username).Return(true, nil)
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, userService := setupCommandUserService(t, tt)
			err := userService.RegisterUser.Handle(ctx, tt.command)
			assertError(t, err, tt.expectedErr)
		})
	}
}

func setupCommandUserService[T any](t *testing.T, tt commandTestCase[T]) (context.Context, *service.Application) {
	t.Helper()
	ctx := context.Background()
	ctxWithAuthUser := auth.NewContextWithUser(ctx, tt.authUser)

	userRepo := domain_mocks.NewMockUserRepository(t)
	userReadModelRepo := domain_mocks.NewMockUserReadModelRepository(t)
	guard := guard_mocks.NewMockGuards(t)

	tt.setupMocks(t, userRepo, guard, &tt.command, tt.authUser)

	userService := service.New(userRepo, userReadModelRepo, guard)

	return ctxWithAuthUser, userService
}
func setupQueryUserService[T any, D any](t *testing.T, tt queryTestCase[T, D]) (context.Context, *service.Application) {
	t.Helper()
	ctx := context.Background()
	ctxWithAuthUser := auth.NewContextWithUser(ctx, tt.authUser)
	userRepo := domain_mocks.NewMockUserRepository(t)
	userReadModelRepo := domain_mocks.NewMockUserReadModelRepository(t)
	guard := guard_mocks.NewMockGuards(t)

	tt.setupMocks(t, userReadModelRepo, guard, tt.query, tt.authUser)

	userService := service.New(userRepo, userReadModelRepo, guard)

	return ctxWithAuthUser, userService
}

func assertError(t *testing.T, err error, expectedErr error) {
	t.Helper()
	if err != nil {
		require.ErrorIs(t, err, expectedErr, expectedErr.Error())
	} else {
		assert.NoError(t, err)
	}
}
