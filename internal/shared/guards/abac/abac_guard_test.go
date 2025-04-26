package abac_test

import (
	"testing"

	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/guards/abac"
	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name string
	input
	expectedErr error
}
type input struct {
	userId   string
	authUser *auth.AuthenticatedUser
}

func TestGuard_CanChangeUsername(t *testing.T) {
	testCases := []testCase{
		{
			name: "user with regular role can change their username",
			input: input{
				userId: "user1",
				authUser: &auth.AuthenticatedUser{
					Id:    "user1",
					Email: "user1@example.com",
					Role:  rbac.Regular,
				},
			},
			expectedErr: nil,
		},
		{
			name: "user with moderator role can change their username",
			input: input{
				userId: "user1",
				authUser: &auth.AuthenticatedUser{
					Id:    "user1",
					Email: "user1@example.com",
					Role:  rbac.Moderator,
				},
			},
			expectedErr: nil,
		},
		{
			name: "user with admin role can change username",
			input: input{
				userId: "user1",
				authUser: &auth.AuthenticatedUser{
					Id:    "user1",
					Email: "user1@example.com",
					Role:  rbac.Admin,
				},
			},
			expectedErr: nil,
		},
		{
			name: "user with regular role cannot change another user's username",
			input: input{
				userId: "user2",
				authUser: &auth.AuthenticatedUser{
					Id:    "user1",
					Email: "user1@example.com",
					Role:  rbac.Regular,
				},
			},
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name: "user with moderator role cannot change another user's username",
			input: input{
				userId: "user2",
				authUser: &auth.AuthenticatedUser{
					Id:    "user1",
					Email: "user1@example.com",
					Role:  rbac.Moderator,
				},
			},
			expectedErr: rbac.ErrUnauthorized,
		},
		{
			name: "user with guest role cannot change another user's username",
			input: input{
				userId: "user1",
				authUser: &auth.AuthenticatedUser{
					Id:    "",
					Email: "",
					Role:  rbac.Guest,
				},
			},
			expectedErr: rbac.ErrUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := abac.New().CanChangeUsername(tc.input.userId, tc.input.authUser)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
