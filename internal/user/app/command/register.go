package command

import (
	"context"
	"time"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/lucsky/cuid"
)

type RegisterUserCommand struct {
	Email    string
	Username string
}

type RegisterUserHandler = shared.CommandHandler[RegisterUserCommand]

type registerUserCommandHandler struct {
	userRepo domain.UserRepository
	guard    rbac.RequestGuard
}

func NewRegisterUserCommandHandler(userRepo domain.UserRepository, guard rbac.RequestGuard) RegisterUserHandler {
	if userRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &registerUserCommandHandler{userRepo: userRepo, guard: guard}
}

func (r *registerUserCommandHandler) Handle(ctx context.Context, cmd RegisterUserCommand) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := r.guard.Authorize(authUser.Role, rbac.CreateAccount); err != nil {
		return err
	}
	user, err := domain.NewUser(
		cuid.New(), cmd.Email, cmd.Username,
		rbac.Regular, time.Now(), time.Now(), nil)
	if err != nil {
		return err
	}
	return r.userRepo.Register(ctx, user)
}
