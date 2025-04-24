package domain

import (
	"fmt"
	"slices"

	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
)

func (u *User) IsModerator() bool {
	return u.role == rbac.Moderator
}

func (u *User) IsAdmin() bool {
	return u.role == rbac.Admin
}

func (u *User) MakeModerator() error {
	if u.role == rbac.Moderator {
		return fmt.Errorf("the user %s is already a moderator", u.username)
	}
	u.role = rbac.Moderator
	return nil
}

func (u *User) IsRegular() bool {
	return u.role == rbac.Regular
}

func (u *User) Role() rbac.UserRole {
	return u.role
}

func (u *User) MakeRegular() error {
	if u.role == rbac.Regular {
		return fmt.Errorf("the user %s is already a regular", u.username)
	}
	u.role = rbac.Regular

	return nil
}

func isValidRole(role rbac.UserRole) bool {
	validRoles := []rbac.UserRole{rbac.Moderator, rbac.Admin, rbac.Regular}

	return slices.Contains(validRoles, role)
}
