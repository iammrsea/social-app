package command

import (
	"context"
	"crypto/rand"
	"errors"
	"time"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type RegisterUserCommand struct {
	Email    string
	Username string
}

type RegisterUserHandler = shared.CommandHandler[RegisterUserCommand]

type registerUserCommandHandler struct {
	userRepo domain.UserRepository
}

func NewRegisterUserCommandHandler(userRepo domain.UserRepository) RegisterUserHandler {
	if userRepo == nil {
		panic("nil user Repository")
	}
	return &registerUserCommandHandler{userRepo: userRepo}
}

func (r *registerUserCommandHandler) Handle(ctx context.Context, cmd RegisterUserCommand) error {
	authUser := auth.GetUserFromCtx(ctx)
	if authUser.IsZero() {
		return errors.New("you are already signed in")
	}

	id := rand.Text()
	user, err := domain.NewUser(id,
		cmd.Email, cmd.Username,
		domain.Regular,
		time.Now(),
		time.Now(),
		nil)
	if err != nil {
		return err
	}

	err = r.userRepo.Register(ctx, user)

	if err != nil {
		return err
	}

	return nil
}
