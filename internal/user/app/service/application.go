package service

import (
	"context"

	"github.com/iammrsea/social-app/internal/user/adapters"
	"github.com/iammrsea/social-app/internal/user/app"
	"github.com/iammrsea/social-app/internal/user/app/commands"
	"github.com/iammrsea/social-app/internal/user/app/queries"
)

func NewApplication(ctx context.Context) app.Application {
	userRepo := adapters.NewMemoryRepository(ctx)

	return app.Application{
		CommandHandlers: app.CommandHandlers{
			RegisterUserHandler:       commands.NewRegisterUserCommandHandler(userRepo),
			RevokeAwardedBagdeHandler: commands.NewRevokeAwardedBadgeCommandHandler(userRepo),
			AwardBadgeHandler:         commands.NewAwardBadgeCommandHandler(userRepo),
			MakeModeratorHandler:      commands.NewMakeModeratorCommandHandler(userRepo),
			ChangeUsernameHandler:     commands.NewChangeUsernameCommandHandler(userRepo),
		},
		QueryHandlers: app.QueryHandlers{
			GetUserHandler:  queries.NewGetUserByIdCommandHandler(userRepo),
			GetUsersHandler: queries.NewGetUsersCommandHandler(userRepo),
		},
	}
}
