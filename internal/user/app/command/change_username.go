package command

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type ChangeUsernameCommand struct {
	Id       string
	Username string
}

type ChangeUsernameHandler = shared.CommandHandler[ChangeUsernameCommand]

type changeUsernameCommandHandler struct {
	userRepo domain.UserRepository
}

func NewChangeUsernameCommandHandler(userRep domain.UserRepository) ChangeUsernameHandler {
	if userRep == nil {
		panic("nil user repository")
	}
	return &changeUsernameCommandHandler{userRepo: userRep}
}

func (c *changeUsernameCommandHandler) Handle(ctx context.Context, cmd ChangeUsernameCommand) error {
	authUser := auth.GetUserFromCtx(ctx)
	if !authUser.IsAuthenticated() {
		return errors.New("unauthorized")
	}
	if authUser.Id != cmd.Id {
		return errors.New("username cannot be changed by proxy")
	}
	err := c.userRepo.ChangeUsername(ctx, cmd.Id, func(user *domain.User) (*domain.User, error) {
		err := user.ChangeUsername(cmd.Username)
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
