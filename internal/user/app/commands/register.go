package commands

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	userDomain "github.com/iammrsea/social-app/internal/user/domain"
)

type RegisterUserCommand struct {
	Email    string
	Username string
}

type RegisterUserHandler = shared.CommandHandler[RegisterUserCommand]

type registerUserCommandHandler struct {
	userRepo userDomain.Repository
}

func NewRegisterUserCommandHandler(userRepo userDomain.Repository) RegisterUserHandler {
	if userRepo == nil {
		panic("nil user Repository")
	}
	return &registerUserCommandHandler{userRepo: userRepo}
}

func (r *registerUserCommandHandler) Handle(ctx context.Context, cmd RegisterUserCommand) error {
	id := rand.Text()
	user, err := userDomain.NewRegularUser(id, cmd.Email, cmd.Username, nil)
	if err != nil {
		return errors.Unwrap(err)
	}

	err = r.userRepo.Register(ctx, user)

	if err != nil {
		return errors.Unwrap(err)
	}

	return nil
}
