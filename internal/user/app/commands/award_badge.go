package commands

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	userDomain "github.com/iammrsea/social-app/internal/user/domain"
)

type AwardBadgeCommand struct {
	Id           string
	Badge        string
	LoggedInUser *userDomain.User
}

type AwardBadgeHandler = shared.CommandHandler[AwardBadgeCommand]

type awardBagdeCommandHandler struct {
	userRepo userDomain.Repository
}

func NewAwardBadgeCommandHandler(userRepo userDomain.Repository) AwardBadgeHandler {
	if userRepo == nil {
		panic("nil user Repository")
	}
	return &awardBagdeCommandHandler{userRepo: userRepo}
}

func (a *awardBagdeCommandHandler) Handle(ctx context.Context, cmd AwardBadgeCommand) error {
	if !cmd.LoggedInUser.IsAdmin() {
		return errors.New("only an admin can award a badge to a user")
	}
	err := a.userRepo.AwardBadge(ctx, cmd.Id, func(user *userDomain.User) (*userDomain.User, error) {
		err := user.AwardBadge(cmd.Badge)
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
