package command

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type MakeModeratorCommand struct {
	Id string
}

type MakeModeratorHandler = shared.CommandHandler[MakeModeratorCommand]

type makeModeratorCommandHandler struct {
	userRepo domain.UserRepository
}

func NewMakeModeratorCommandHandler(userRepo domain.UserRepository) MakeModeratorHandler {
	if userRepo == nil {
		panic("nil user Repository")
	}
	return &makeModeratorCommandHandler{userRepo: userRepo}
}

func (r *makeModeratorCommandHandler) Handle(ctx context.Context, cmd MakeModeratorCommand) error {
	authUser := auth.GetUserFromCtx(ctx)
	if !authUser.IsAuthenticated() {
		return errors.New("unauthorized")
	}
	if authUser.Role != domain.Admin {
		return errors.New("only an admin can change user role")
	}
	err := r.userRepo.MakeModerator(ctx, cmd.Id, func(user *domain.User) (*domain.User, error) {
		err := user.MakeModerator()
		if err != nil {
			return user, err
		}
		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}
