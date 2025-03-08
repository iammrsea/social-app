package domain

import (
	"fmt"
	"slices"
)

type UserRole string

const (
	admin     UserRole = "admin"
	moderator UserRole = "moderator"
	regular   UserRole = "regular"
)

func (u *User) IsModerator() bool {
	return u.role == moderator
}

func (u *User) IsAdmin() bool {
	return u.role == admin
}

func (u *User) MakeModerator() error {
	if u.role == moderator {
		return fmt.Errorf("the user %s is already a moderator", u.username)
	}
	u.role = moderator
	return nil
}

func (u *User) IsRegular() bool {
	return u.role == regular
}

func (u *User) Role() string {
	return string(u.role)
}

func (u *User) MakeRegular() error {
	if u.role == regular {
		return fmt.Errorf("the user %s is already a regular", u.username)
	}
	u.role = regular

	return nil
}

func isValidRole(role UserRole) bool {
	validRoles := []UserRole{moderator, admin, regular}

	return slices.Contains(validRoles, role)
}
