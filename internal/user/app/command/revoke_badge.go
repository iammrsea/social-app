package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type RevokeAwardedBadgeCommand struct {
	Id    string
	Badge string
}

type RevokeAwardedBagdeHandler = shared.CommandHandler[RevokeAwardedBadgeCommand]

type revokeAwardedBagdeCommandHandler struct {
	userRepo domain.UserRepository
	guard    rbac.RequestGuard
}

func NewRevokeAwardedBadgeCommandHandler(userRepo domain.UserRepository, guard rbac.RequestGuard) RevokeAwardedBagdeHandler {
	if userRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &revokeAwardedBagdeCommandHandler{userRepo: userRepo, guard: guard}
}

func (r *revokeAwardedBagdeCommandHandler) Handle(ctx context.Context, cmd RevokeAwardedBadgeCommand) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := r.guard.Authorize(authUser.Role, rbac.RevokeBadge); err != nil {
		return err
	}
	return r.userRepo.RevokeAwardedBadge(ctx, cmd.Id, func(user *domain.User) error {
		return user.RevokeAwardedBadge(cmd.Badge)
	})
}
