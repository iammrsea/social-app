package rbac_test

import (
	"testing"

	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	userRole    rbac.UserRole
	permission  rbac.Permission
	expectedErr error
}

func TestPermision_BanUser(t *testing.T) {
	t.Parallel()
	testCases := []testCase{
		{
			name:        "user with admin role can ban user",
			userRole:    rbac.Admin,
			permission:  rbac.BanUser,
			expectedErr: nil,
		},
		{
			name:        "user with moderator role can ban user",
			userRole:    rbac.Moderator,
			permission:  rbac.BanUser,
			expectedErr: nil,
		},
		{
			name:        "user with regular role cannot ban user",
			userRole:    rbac.Regular,
			permission:  rbac.BanUser,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with guest role cannot ban user",
			userRole:    rbac.Guest,
			permission:  rbac.BanUser,
			expectedErr: rbac.ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := rbac.New().Authorize(tc.userRole, tc.permission)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestPermission_UnbanUser(t *testing.T) {
	t.Parallel()
	testCases := []testCase{
		{
			name:        "user with admin role can unban user",
			userRole:    rbac.Admin,
			permission:  rbac.UnbanUser,
			expectedErr: nil,
		},
		{
			name:        "user with moderator role can unban user",
			userRole:    rbac.Moderator,
			permission:  rbac.UnbanUser,
			expectedErr: nil,
		},
		{
			name:        "user with regular role cannot unban user",
			userRole:    rbac.Regular,
			permission:  rbac.UnbanUser,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with guest role cannot unban user",
			userRole:    rbac.Guest,
			permission:  rbac.UnbanUser,
			expectedErr: rbac.ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := rbac.New().Authorize(tc.userRole, tc.permission)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestPermission_ViewUser(t *testing.T) {
	t.Parallel()
	testCases := []testCase{
		{
			name:        "user with admin role can view user",
			userRole:    rbac.Admin,
			permission:  rbac.ViewUser,
			expectedErr: nil,
		},
		{
			name:        "user with moderator role can view user",
			userRole:    rbac.Moderator,
			permission:  rbac.ViewUser,
			expectedErr: nil,
		},
		{
			name:        "user with regular role can view user",
			userRole:    rbac.Regular,
			permission:  rbac.ViewUser,
			expectedErr: nil,
		},
		{
			name:        "user with guest role cannot view user",
			userRole:    rbac.Guest,
			permission:  rbac.ViewUser,
			expectedErr: rbac.ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := rbac.New().Authorize(tc.userRole, tc.permission)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestPermission_DeleteUser(t *testing.T) {
	t.Parallel()
	testCases := []testCase{
		{
			name:        "user with admin role can delete user",
			userRole:    rbac.Admin,
			permission:  rbac.DeleteUser,
			expectedErr: nil,
		},
		{
			name:        "user with moderator role cannot delete user",
			userRole:    rbac.Moderator,
			permission:  rbac.DeleteUser,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with regular role cannot delete user",
			userRole:    rbac.Regular,
			permission:  rbac.DeleteUser,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with guest role cannot delete user",
			userRole:    rbac.Guest,
			permission:  rbac.DeleteUser,
			expectedErr: rbac.ErrUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := rbac.New().Authorize(tc.userRole, tc.permission)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestPermission_AwardBadge(t *testing.T) {
	t.Parallel()
	testCases := []testCase{
		{
			name:        "user with admin role can award badge",
			userRole:    rbac.Admin,
			permission:  rbac.AwardBadge,
			expectedErr: nil,
		},
		{
			name:        "user with moderator role cannot award badge",
			userRole:    rbac.Moderator,
			permission:  rbac.AwardBadge,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with regular role cannot award badge",
			userRole:    rbac.Regular,
			permission:  rbac.AwardBadge,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with guest role cannot award badge",
			userRole:    rbac.Guest,
			permission:  rbac.AwardBadge,
			expectedErr: rbac.ErrUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := rbac.New().Authorize(tc.userRole, tc.permission)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestPermission_RevokeBadge(t *testing.T) {
	t.Parallel()
	testCases := []testCase{
		{
			name:        "user with admin role can revoke badge",
			userRole:    rbac.Admin,
			permission:  rbac.RevokeBadge,
			expectedErr: nil,
		},
		{
			name:        "user with moderator role cannot revoke badge",
			userRole:    rbac.Moderator,
			permission:  rbac.RevokeBadge,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with regular role cannot revoke badge",
			userRole:    rbac.Regular,
			permission:  rbac.RevokeBadge,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with guest role cannot revoke badge",
			userRole:    rbac.Guest,
			permission:  rbac.RevokeBadge,
			expectedErr: rbac.ErrUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := rbac.New().Authorize(tc.userRole, tc.permission)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestPermission_MakeModerator(t *testing.T) {
	t.Parallel()
	testCases := []testCase{
		{
			name:        "user with admin role can make moderator",
			userRole:    rbac.Admin,
			permission:  rbac.MakeModerator,
			expectedErr: nil,
		},
		{
			name:        "user with moderator role cannot make moderator",
			userRole:    rbac.Moderator,
			permission:  rbac.MakeModerator,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with regular role cannot make moderator",
			userRole:    rbac.Regular,
			permission:  rbac.MakeModerator,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with guest role cannot make moderator",
			userRole:    rbac.Guest,
			permission:  rbac.MakeModerator,
			expectedErr: rbac.ErrUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := rbac.New().Authorize(tc.userRole, tc.permission)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestPermission_MakeRegular(t *testing.T) {
	t.Parallel()
	testCases := []testCase{
		{
			name:        "user with admin role can make regular",
			userRole:    rbac.Admin,
			permission:  rbac.MakeRegular,
			expectedErr: nil,
		},
		{
			name:        "user with moderator role cannot make regular",
			userRole:    rbac.Moderator,
			permission:  rbac.MakeRegular,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with regular role cannot make regular",
			userRole:    rbac.Regular,
			permission:  rbac.MakeRegular,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with guest role cannot make regular",
			userRole:    rbac.Guest,
			permission:  rbac.MakeRegular,
			expectedErr: rbac.ErrUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := rbac.New().Authorize(tc.userRole, tc.permission)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestPermission_CreateAccount(t *testing.T) {
	t.Parallel()
	testCases := []testCase{
		{
			name:        "user with admin role can create account",
			userRole:    rbac.Admin,
			permission:  rbac.CreateAccount,
			expectedErr: nil,
		},
		{
			name:        "user with moderator role cannot create account",
			userRole:    rbac.Moderator,
			permission:  rbac.CreateAccount,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with regular role cannot create account",
			userRole:    rbac.Regular,
			permission:  rbac.CreateAccount,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with guest role can create account",
			userRole:    rbac.Guest,
			permission:  rbac.CreateAccount,
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := rbac.New().Authorize(tc.userRole, tc.permission)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestPermission_ListUsers(t *testing.T) {
	t.Parallel()
	testCases := []testCase{
		{
			name:        "user with admin role can list users",
			userRole:    rbac.Admin,
			permission:  rbac.ListUsers,
			expectedErr: nil,
		},
		{
			name:        "user with moderator role can list users",
			userRole:    rbac.Moderator,
			permission:  rbac.ListUsers,
			expectedErr: nil,
		},
		{
			name:        "user with regular role cannot list users",
			userRole:    rbac.Regular,
			permission:  rbac.ListUsers,
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name:        "user with guest role cannot list users",
			userRole:    rbac.Guest,
			permission:  rbac.ListUsers,
			expectedErr: rbac.ErrUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := rbac.New().Authorize(tc.userRole, tc.permission)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
