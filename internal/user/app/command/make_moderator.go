package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/guards"
	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type MakeModerator struct {
	Id string
}

type MakeModeratorHandler = shared.CommandHandler[MakeModerator]

type makeModeratorHandler struct {
	userRepo domain.UserRepository
	guard    guards.Guards
}

func NewMakeModeratorHandler(userRepo domain.UserRepository, guard guards.Guards) MakeModeratorHandler {
	if userRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &makeModeratorHandler{userRepo: userRepo, guard: guard}
}

func (r *makeModeratorHandler) Handle(ctx context.Context, cmd MakeModerator) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := r.guard.Authorize(authUser.Role, rbac.MakeModerator); err != nil {
		return err
	}
	return r.userRepo.MakeModerator(ctx, cmd.Id, func(user *domain.User) error {
		return user.MakeModerator()
	})
}
