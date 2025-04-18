package domain

import (
	"fmt"
	"slices"
)

type UserRole string

const (
	Admin     UserRole = "ADMIN"
	Moderator UserRole = "MODERATOR"
	Regular   UserRole = "REGULAR"
)

func (r UserRole) String() string {
	return string(r)
}

func (u *User) IsModerator() bool {
	return u.role == Moderator
}

func (u *User) IsAdmin() bool {
	return u.role == Admin
}

func (u *User) MakeModerator() error {
	if u.role == Moderator {
		return fmt.Errorf("the user %s is already a moderator", u.username)
	}
	u.role = Moderator
	return nil
}

func (u *User) IsRegular() bool {
	return u.role == Regular
}

func (u *User) Role() UserRole {
	return u.role
}

func (u *User) MakeRegular() error {
	if u.role == Regular {
		return fmt.Errorf("the user %s is already a regular", u.username)
	}
	u.role = Regular

	return nil
}

func isValidRole(role UserRole) bool {
	validRoles := []UserRole{Moderator, Admin, Regular}

	return slices.Contains(validRoles, role)
}
