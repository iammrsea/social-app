package guards

import (
	"github.com/iammrsea/social-app/internal/shared/guards/abac"
	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
)

type Guards interface {
	// RBAC
	rbac.Guard
	// ABAC
	abac.Guard
}

type guards struct {
	*rbac.RoleBasedGuard
	*abac.AttributeBasedGuard
}

func New() Guards {
	return &guards{
		RoleBasedGuard:      rbac.New(),
		AttributeBasedGuard: abac.New(),
	}
}

func LoadUsers(guard Guards) {
	guard.Authorize(rbac.Admin, rbac.Permission("read"))
	// guard.CanChangeUsername()
}

func Run() {
	guard := New()

	LoadUsers(guard)
}
