package rbac

//rbac => Role-Based Access Control

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Guard interface {
	Authorize(role UserRole, perm Permission) error
}

type policy interface {
	IsAllowed(role UserRole, perm Permission) bool
}

type RoleBasedGuard struct {
	policy
}

func New() *RoleBasedGuard {
	return &RoleBasedGuard{policy: NewPolicy()}
}

func (rg *RoleBasedGuard) Authorize(role UserRole, perm Permission) error {
	if !rg.policy.IsAllowed(role, perm) {
		return ErrUnauthorized
	}
	return nil
}
