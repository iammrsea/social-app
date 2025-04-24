package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/guards"
	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type RevokeAwardedBadge struct {
	Id    string
	Badge string
}

type RevokeAwardedBadgeHandler = shared.CommandHandler[RevokeAwardedBadge]

type revokeAwardedBagdeHandler struct {
	userRepo domain.UserRepository
	guard    guards.Guards
}

func NewRevokeAwardedBadgeHandler(userRepo domain.UserRepository, guard guards.Guards) RevokeAwardedBadgeHandler {
	if userRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &revokeAwardedBagdeHandler{userRepo: userRepo, guard: guard}
}

func (r *revokeAwardedBagdeHandler) Handle(ctx context.Context, cmd RevokeAwardedBadge) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := r.guard.Authorize(authUser.Role, rbac.RevokeBadge); err != nil {
		return err
	}
	return r.userRepo.RevokeAwardedBadge(ctx, cmd.Id, func(user *domain.User) error {
		return user.RevokeAwardedBadge(cmd.Badge)
	})
}
