package command

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/user/abac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type ChangeUsername struct {
	Id       string
	Username string
}

type ChangeUsernameHandler = shared.CommandHandler[ChangeUsername]

type changeUsernameHandler struct {
	userRepo domain.UserRepository
}

func NewChangeUsernameHandler(userRep domain.UserRepository) ChangeUsernameHandler {
	if userRep == nil {
		panic("nil user repository")
	}
	return &changeUsernameHandler{userRepo: userRep}
}

func (c *changeUsernameHandler) Handle(ctx context.Context, cmd ChangeUsername) error {
	authUser := auth.GetUserFromCtx(ctx)
	if err := abac.CanChangeUsername(cmd.Id, authUser); err != nil {
		return err
	}
	return c.userRepo.ChangeUsername(ctx, cmd.Id, func(user *domain.User) error {
		return user.ChangeUsername(cmd.Username)
	})
}
