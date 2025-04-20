package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type AwardBadge struct {
	Id    string
	Badge string
}

type AwardBadgeHandler = shared.CommandHandler[AwardBadge]

type awardBadgeHandler struct {
	userRepo domain.UserRepository
	guard    rbac.RequestGuard
}

func NewAwardBadgeHandler(userRepo domain.UserRepository, guard rbac.RequestGuard) AwardBadgeHandler {
	if userRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &awardBadgeHandler{userRepo: userRepo, guard: guard}
}

func (a *awardBadgeHandler) Handle(ctx context.Context, cmd AwardBadge) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := a.guard.Authorize(authUser.Role, rbac.AwardBadge); err != nil {
		return err
	}
	return a.userRepo.AwardBadge(ctx, cmd.Id, func(user *domain.User) error {
		return user.AwardBadge(cmd.Badge)
	})
}
