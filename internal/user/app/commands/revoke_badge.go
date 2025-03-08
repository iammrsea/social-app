package commands

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	userDomain "github.com/iammrsea/social-app/internal/user/domain"
)

type RevokeAwardedBadgeCommand struct {
	Id           string
	Badge        string
	LoggedInUser *userDomain.User
}

type RevokeAwardedBagdeHandler = shared.CommandHandler[RevokeAwardedBadgeCommand]

type revokeAwardedBagdeCommandHandler struct {
	userRepo userDomain.Repository
}

func NewRevokeAwardedBadgeCommandHandler(userRepo userDomain.Repository) RevokeAwardedBagdeHandler {
	if userRepo == nil {
		panic("nil user Repository")
	}
	return &revokeAwardedBagdeCommandHandler{userRepo: userRepo}
}

func (r *revokeAwardedBagdeCommandHandler) Handle(ctx context.Context, cmd RevokeAwardedBadgeCommand) error {
	if !cmd.LoggedInUser.IsAdmin() {
		return errors.New("only an admin can revoke a badge awarded to a user")
	}
	err := r.userRepo.RevokeAwardedBadge(ctx, cmd.Id, func(user *userDomain.User) (*userDomain.User, error) {
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
