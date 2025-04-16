package command

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type RevokeAwardedBadgeCommand struct {
	Id           string
	Badge        string
	LoggedInUser *domain.User
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
	if !cmd.LoggedInUser.IsAdmin() {
		return errors.New("only an admin can revoke a badge awarded to a user")
	}
	err := r.userRepo.RevokeAwardedBadge(ctx, cmd.Id, func(user *domain.User) (*domain.User, error) {
		err := user.RevokeAwardedBadge(cmd.Badge)
		if err != nil {
			return nil, errors.Unwrap(err)
		}
		return user, nil
	})

	if err != nil {
		return errors.Unwrap(err)
	}

	return nil
}
