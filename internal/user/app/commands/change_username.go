package commands

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	userDomain "github.com/iammrsea/social-app/internal/user/domain"
)

type ChangeUsernameCommand struct {
	Id           string
	Username     string
	LoggedInUser *userDomain.User
}

type ChangeUsernameHandler = shared.CommandHandler[ChangeUsernameCommand]

type changeUsernameCommandHandler struct {
	userRepo userDomain.Repository
}

func NewChangeUsernameCommandHandler(userRep userDomain.Repository) ChangeUsernameHandler {
	if userRep == nil {
		panic("nil user repository")
	}
	return &changeUsernameCommandHandler{userRepo: userRep}
}

func (c *changeUsernameCommandHandler) Handle(ctx context.Context, cmd ChangeUsernameCommand) error {
	if cmd.LoggedInUser.Id() != cmd.Id {
		return errors.New("username cannot be changed by proxy")
	}
	err := c.userRepo.ChangeUsername(ctx, cmd.Id, func(user *userDomain.User) (*userDomain.User, error) {
		err := user.ChangeUsername(cmd.Username)
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
