package command

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type AwardBadgeCommand struct {
	Id    string
	Badge string
}

type AwardBadgeHandler = shared.CommandHandler[AwardBadgeCommand]

type awardBagdeCommandHandler struct {
	userRepo domain.UserRepository
}

func NewAwardBadgeCommandHandler(userRepo domain.UserRepository) AwardBadgeHandler {
	if userRepo == nil {
		panic("nil user Repository")
	}
	return &awardBagdeCommandHandler{userRepo: userRepo}
}

func (a *awardBagdeCommandHandler) Handle(ctx context.Context, cmd AwardBadgeCommand) error {
	authUser := auth.GetUserFromCtx(ctx)
	if !authUser.IsAuthenticated() {
		return errors.New("unauthorized")
	}
	if authUser.Role != domain.Admin {
		return errors.New("only an admin can award a badge to a user")
	}
	err := a.userRepo.AwardBadge(ctx, cmd.Id, func(user *domain.User) (*domain.User, error) {
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
