package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type MakeModeratorCommand struct {
	Id string
}

type MakeModeratorHandler = shared.CommandHandler[MakeModeratorCommand]

type makeModeratorCommandHandler struct {
	userRepo domain.UserRepository
	guard    rbac.RequestGuard
}

func NewMakeModeratorCommandHandler(userRepo domain.UserRepository, guard rbac.RequestGuard) MakeModeratorHandler {
	if userRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &makeModeratorCommandHandler{userRepo: userRepo, guard: guard}
}

func (r *makeModeratorCommandHandler) Handle(ctx context.Context, cmd MakeModeratorCommand) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := r.guard.Authorize(authUser.Role, rbac.MakeModerator); err != nil {
		return err
	}
	err := r.userRepo.MakeModerator(ctx, cmd.Id, func(user *domain.User) error {
		return user.MakeModerator()
	})

	return err
}
