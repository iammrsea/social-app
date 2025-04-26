package abac

import (
	"github.com/iammrsea/social-app/internal/shared/auth"
)

//ABAC => Attribute-Based Access Control

type Guard interface {
	CanChangeUsername(userId string, authUser *auth.AuthenticatedUser) error
}

type AttributeBasedGuard struct{}

func New() *AttributeBasedGuard {
	return &AttributeBasedGuard{}
}
