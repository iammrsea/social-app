package commands

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	userDomain "github.com/iammrsea/social-app/internal/user/domain"
)

type MakeModeratorCommand struct {
	Id string
}

type MakeModeratorHandler = shared.CommandHandler[MakeModeratorCommand]

type makeModeratorCommandHandler struct {
	userRepo userDomain.Repository
}

func NewMakeModeratorCommandHandler(userRepo userDomain.Repository) MakeModeratorHandler {
	if userRepo == nil {
		panic("nil user Repository")
	}
	return &makeModeratorCommandHandler{userRepo: userRepo}
}

func (r *makeModeratorCommandHandler) Handle(ctx context.Context, cmd MakeModeratorCommand) error {
	err := r.userRepo.MakeModerator(ctx, cmd.Id, func(user *userDomain.User) (*userDomain.User, error) {
		err := user.MakeModerator()
		if err != nil {
			return user, err
		}
		return nil, nil
	})

	if err != nil {
		return errors.Unwrap(err)
	}

	return nil
}
