package service

import (
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/app"
	"github.com/iammrsea/social-app/internal/user/app/command"
	"github.com/iammrsea/social-app/internal/user/app/query"
	"github.com/iammrsea/social-app/internal/user/domain"
)

// Constructor of the user application layer
func NewUserService(userRepo domain.UserRepository, userReadModelRepo domain.UserReadModelRepository, guard rbac.RequestGuard) *app.Application {
	return &app.Application{
		CommandHandler: app.CommandHandler{
			RegisterUserHandler:       command.NewRegisterUserCommandHandler(userRepo, guard),
			RevokeAwardedBagdeHandler: command.NewRevokeAwardedBadgeCommandHandler(userRepo, guard),
			AwardBadgeHandler:         command.NewAwardBadgeCommandHandler(userRepo, guard),
			MakeModeratorHandler:      command.NewMakeModeratorCommandHandler(userRepo, guard),
			ChangeUsernameHandler:     command.NewChangeUsernameCommandHandler(userRepo),
			BanUserHandler:            command.NewBanUserHandler(userRepo, guard),
		},
		QueryHandler: app.QueryHandler{
			GetUserByIdHandler:    query.NewGetUserByIdCommandHandler(userReadModelRepo, guard),
			GetUsersHandler:       query.NewGetUsersCommandHandler(userReadModelRepo, guard),
			GetUserByEmailHandler: query.NewGetUserByEmailCommandHandler(userReadModelRepo, guard),
		},
	}
}
