package rbac

import (
	"slices"
)

type Policy struct {
	rules map[UserRole][]Permission
}

func NewPolicy() *Policy {
	return &Policy{
		rules: map[UserRole][]Permission{
			Regular:   {ViewUser},
			Admin:     {ViewUser},
			Moderator: {ViewUser, ListUsers},
			Guest:     {},
		},
	}
}

func (p *Policy) IsAllowed(role UserRole, perm Permission) bool {
	if role == Admin {
		return true
	}
	perms, ok := p.rules[role]
	if !ok {
		return false
	}
	return slices.Contains(perms, perm)
}
