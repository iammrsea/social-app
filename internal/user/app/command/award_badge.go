package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type AwardBadgeCommand struct {
	Id    string
	Badge string
}

type AwardBadgeHandler = shared.CommandHandler[AwardBadgeCommand]

type awardBagdeCommandHandler struct {
	userRepo domain.UserRepository
	guard    rbac.RequestGuard
}

func NewAwardBadgeCommandHandler(userRepo domain.UserRepository, guard rbac.RequestGuard) AwardBadgeHandler {
	if userRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &awardBagdeCommandHandler{userRepo: userRepo, guard: guard}
}

func (a *awardBagdeCommandHandler) Handle(ctx context.Context, cmd AwardBadgeCommand) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := a.guard.Authorize(authUser.Role, rbac.AwardBadge); err != nil {
		return err
	}
	return a.userRepo.AwardBadge(ctx, cmd.Id, func(user *domain.User) error {
		return user.AwardBadge(cmd.Badge)
	})
}
