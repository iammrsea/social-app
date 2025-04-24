package abac

import (
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
)

func (g *AttributeBasedGuard) CanChangeUsername(userId string, authUser *auth.AuthenticatedUser) error {
	switch authUser.Role {
	case rbac.Admin:
		return nil
	case rbac.Regular, rbac.Moderator:
		if userId == authUser.Id {
			return nil
		} else {
			return rbac.ErrUnauthorized
		}
	default:
		return rbac.ErrUnauthorized
	}
}
