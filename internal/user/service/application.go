package service

import (
	"context"

	"github.com/iammrsea/social-app/internal/user/app"
	"github.com/iammrsea/social-app/internal/user/app/command"
	"github.com/iammrsea/social-app/internal/user/app/query"
	"github.com/iammrsea/social-app/internal/user/infra/db/memory"
)

// Constructor of the user application layer
func NewUserService(ctx context.Context) *app.Application {
	userRepo := memory.NewUserRepository(ctx)

	return &app.Application{
		CommandHandler: app.CommandHandler{
			RegisterUserHandler:       command.NewRegisterUserCommandHandler(userRepo),
			RevokeAwardedBagdeHandler: command.NewRevokeAwardedBadgeCommandHandler(userRepo),
			AwardBadgeHandler:         command.NewAwardBadgeCommandHandler(userRepo),
			MakeModeratorHandler:      command.NewMakeModeratorCommandHandler(userRepo),
			ChangeUsernameHandler:     command.NewChangeUsernameCommandHandler(userRepo),
		},
		QueryHandler: app.QueryHandler{
			GetUserByIdHandler:    query.NewGetUserByIdCommandHandler(userRepo),
			GetUsersHandler:       query.NewGetUsersCommandHandler(userRepo),
			GetUserByEmailHandler: query.NewGetUserByEmailCommandHandler(userRepo),
		},
	}
}
