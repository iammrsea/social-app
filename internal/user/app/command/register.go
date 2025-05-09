package command

import (
	"context"
	"errors"
	"time"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/custom_errors"
	"github.com/iammrsea/social-app/internal/shared/guards"
	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/lucsky/cuid"
)

type RegisterUser struct {
	Email    string
	Username string
}

type RegisterUserHandler = shared.CommandHandler[RegisterUser]

type registerUserHandler struct {
	userRepo domain.UserRepository
	guard    guards.Guards
}

func NewRegisterUserHandler(userRepo domain.UserRepository, guard guards.Guards) RegisterUserHandler {
	if userRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &registerUserHandler{userRepo: userRepo, guard: guard}
}

func (r *registerUserHandler) Handle(ctx context.Context, cmd RegisterUser) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := r.guard.Authorize(authUser.Role, rbac.CreateAccount); err != nil {
		return err
	}

	userExists, err := r.userRepo.UserExists(ctx, cmd.Email, cmd.Username)

	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return custom_errors.ErrInternalServerError
		}
	}
	if userExists {
		return domain.ErrEmailOrUsernameAlreadyExists
	}

	user, err := domain.NewUser(
		cuid.New(), cmd.Email, cmd.Username,
		rbac.Regular, time.Now(), time.Now(), nil, nil)
	if err != nil {
		return err
	}
	return r.userRepo.Register(ctx, user)
}
