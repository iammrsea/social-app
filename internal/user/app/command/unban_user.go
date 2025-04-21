package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type UnbanUser struct {
	Id string
}

type UnbanUserHandler = shared.CommandHandler[UnbanUser]

type unbanUserHandler struct {
	userRepo domain.UserRepository
	guard    rbac.RequestGuard
}

func NewUnbanUserHandler(userRepo domain.UserRepository, guard rbac.RequestGuard) UnbanUserHandler {
	if userRepo == nil || guard == nil {
		panic("nil user Repository or guard")
	}
	return &unbanUserHandler{userRepo: userRepo, guard: guard}
}

func (a *unbanUserHandler) Handle(ctx context.Context, cmd UnbanUser) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := a.guard.Authorize(authUser.Role, rbac.UnbanUser); err != nil {
		return err
	}
	return a.userRepo.UnbanUser(ctx, cmd.Id, func(user *domain.User) error {
		return user.UnBan()
	})
}
