package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type BanUserCommand struct {
	Id             string
	Reason         string
	IsIndefinitely bool
	Timeline       *domain.BanTimeline
}

type BanUserHandler = shared.CommandHandler[BanUserCommand]

type banUserHandler struct {
	userRepo domain.UserRepository
	guard    rbac.RequestGuard
}

func NewBanUserHandler(userRepo domain.UserRepository, guard rbac.RequestGuard) BanUserHandler {
	if userRepo == nil || guard == nil {
		panic("nil user Repository or guard")
	}
	return &banUserHandler{userRepo: userRepo, guard: guard}
}

func (a *banUserHandler) Handle(ctx context.Context, cmd BanUserCommand) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := a.guard.Authorize(authUser.Role, rbac.BanUser); err != nil {
		return err
	}
	return a.userRepo.BanUser(ctx, cmd.Id, func(user *domain.User) error {
		return user.Ban(cmd.Reason, cmd.IsIndefinitely, cmd.Timeline)
	})
}
