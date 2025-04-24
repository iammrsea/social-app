package command

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/custom_errors"
	"github.com/iammrsea/social-app/internal/shared/guards"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type ChangeUsername struct {
	Id       string
	Username string
}

type ChangeUsernameHandler = shared.CommandHandler[ChangeUsername]

type changeUsernameHandler struct {
	userRepo domain.UserRepository
	guard    guards.Guards
}

func NewChangeUsernameHandler(userRepo domain.UserRepository, guard guards.Guards) ChangeUsernameHandler {
	if userRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &changeUsernameHandler{userRepo: userRepo, guard: guard}
}

func (c *changeUsernameHandler) Handle(ctx context.Context, cmd ChangeUsername) error {
	authUser := auth.GetUserFromCtx(ctx)
	userExists, err := c.userRepo.UserExists(ctx, "", cmd.Username)

	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return custom_errors.ErrInternalServerError
		}
	}
	if userExists {
		return domain.ErrEmailOrUsernameAlreadyExists
	}
	if err := c.guard.CanChangeUsername(cmd.Id, authUser); err != nil {
		return err
	}
	return c.userRepo.ChangeUsername(ctx, cmd.Id, func(user *domain.User) error {
		return user.ChangeUsername(cmd.Username)
	})
}
