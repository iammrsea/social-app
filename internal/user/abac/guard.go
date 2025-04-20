package abac

//ABAC => Attribute-Based Access Control

import (
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
)

func CanChangeUsername(userId string, authUser *auth.AuthenticatedUser) error {
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
