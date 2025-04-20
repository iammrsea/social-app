package rbac

import "errors"

var ErrUnauthorized = errors.New("unauthorized")

type RequestGuard interface {
	Authorize(role UserRole, perm Permission) error
}

type policy interface {
	IsAllowed(role UserRole, perm Permission) bool
}

type requestGuard struct {
	policy policy
}

func NewRequestGuard(p policy) RequestGuard {
	return &requestGuard{policy: p}
}

func (rg *requestGuard) Authorize(role UserRole, perm Permission) error {
	if !rg.policy.IsAllowed(role, perm) {
		return ErrUnauthorized
	}
	return nil
}
