package command

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type RevokeAwardedBadgeCommand struct {
	Id    string
	Badge string
}

type RevokeAwardedBagdeHandler = shared.CommandHandler[RevokeAwardedBadgeCommand]

type revokeAwardedBagdeCommandHandler struct {
	userRepo domain.UserRepository
}

func NewRevokeAwardedBadgeCommandHandler(userRepo domain.UserRepository) RevokeAwardedBagdeHandler {
	if userRepo == nil {
		panic("nil user Repository")
	}
	return &revokeAwardedBagdeCommandHandler{userRepo: userRepo}
}

func (r *revokeAwardedBagdeCommandHandler) Handle(ctx context.Context, cmd RevokeAwardedBadgeCommand) error {
	authUser := auth.GetUserFromCtx(ctx)
	if !authUser.IsAuthenticated() {
		return errors.New("unauthorized")
	}

	if authUser.Role != domain.Admin {
		return errors.New("only an admin can revoke a badge awarded to a user")
	}
	err := r.userRepo.RevokeAwardedBadge(ctx, cmd.Id, func(user *domain.User) (*domain.User, error) {
		err := user.RevokeAwardedBadge(cmd.Badge)
		if err != nil {
			return nil, err
		}
		return user, nil
	})

	if err != nil {
		return err
	}

	return nil
}
