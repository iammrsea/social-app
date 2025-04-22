package service

import (
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/app/command"
	"github.com/iammrsea/social-app/internal/user/app/query"
	"github.com/iammrsea/social-app/internal/user/domain"
)

// Constructor of the user application layer
func New(userRepo domain.UserRepository, userReadModelRepo domain.UserReadModelRepository, guard rbac.RequestGuard) *Application {
	return &Application{
		CommandHandler: CommandHandler{
			RegisterUser:       command.NewRegisterUserHandler(userRepo, guard),
			RevokeAwardedBadge: command.NewRevokeAwardedBadgeHandler(userRepo, guard),
			AwardBadge:         command.NewAwardBadgeHandler(userRepo, guard),
			MakeModerator:      command.NewMakeModeratorHandler(userRepo, guard),
			ChangeUsername:     command.NewChangeUsernameHandler(userRepo),
			BanUser:            command.NewBanUserHandler(userRepo, guard),
			UnbanUser:          command.NewUnbanUserHandler(userRepo, guard),
		},
		QueryHandler: QueryHandler{
			GetUserById:    query.NewGetUserByIdHandler(userReadModelRepo, guard),
			GetUsers:       query.NewGetUsersHandler(userReadModelRepo, guard),
			GetUserByEmail: query.NewGetUserByEmailHandler(userReadModelRepo, guard),
		},
	}
}
