package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/user/abac"
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
	if err := abac.CanChangeUsername(cmd.Id, authUser); err != nil {
		return err
	}
	return c.userRepo.ChangeUsername(ctx, cmd.Id, func(user *domain.User) error {
		return user.ChangeUsername(cmd.Username)
	})
}
