package service

import (
	"github.com/iammrsea/social-app/internal/user/app"
	"github.com/iammrsea/social-app/internal/user/app/command"
	"github.com/iammrsea/social-app/internal/user/app/query"
	"github.com/iammrsea/social-app/internal/user/domain"
)

// Constructor of the user application layer
func NewUserService(userRepo domain.UserRepository, userReadModelRepo domain.UserReadModelRepository) *app.Application {
	return &app.Application{
		CommandHandler: app.CommandHandler{
			RegisterUserHandler:       command.NewRegisterUserCommandHandler(userRepo),
			RevokeAwardedBagdeHandler: command.NewRevokeAwardedBadgeCommandHandler(userRepo),
			AwardBadgeHandler:         command.NewAwardBadgeCommandHandler(userRepo),
			MakeModeratorHandler:      command.NewMakeModeratorCommandHandler(userRepo),
			ChangeUsernameHandler:     command.NewChangeUsernameCommandHandler(userRepo),
		},
		QueryHandler: app.QueryHandler{
			GetUserByIdHandler:    query.NewGetUserByIdCommandHandler(userReadModelRepo),
			GetUsersHandler:       query.NewGetUsersCommandHandler(userReadModelRepo),
			GetUserByEmailHandler: query.NewGetUserByEmailCommandHandler(userReadModelRepo),
		},
	}
}
