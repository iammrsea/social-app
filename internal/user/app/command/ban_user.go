package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/guards"
	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type BanUser struct {
	Id             string
	Reason         string
	IsIndefinitely bool
	Timeline       *domain.BanTimeline
}

type BanUserHandler = shared.CommandHandler[BanUser]

type banUserHandler struct {
	userRepo domain.UserRepository
	guard    guards.Guards
}

func NewBanUserHandler(userRepo domain.UserRepository, guard guards.Guards) BanUserHandler {
	if userRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &banUserHandler{userRepo: userRepo, guard: guard}
}

func (a *banUserHandler) Handle(ctx context.Context, cmd BanUser) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := a.guard.Authorize(authUser.Role, rbac.BanUser); err != nil {
		return err
	}
	return a.userRepo.BanUser(ctx, cmd.Id, func(user *domain.User) error {
		return user.Ban(cmd.Reason, cmd.IsIndefinitely, cmd.Timeline)
	})
}
