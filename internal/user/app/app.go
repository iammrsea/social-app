package app

import (
	"github.com/iammrsea/social-app/internal/user/app/command"
	"github.com/iammrsea/social-app/internal/user/app/query"
)

type Application struct {
	CommandHandler
	QueryHandler
}

type CommandHandler struct {
	RegisterUser       command.RegisterUserHandler
	RevokeAwardedBadge command.RevokeAwardedBadgeHandler
	AwardBadge         command.AwardBadgeHandler
	MakeModerator      command.MakeModeratorHandler
	ChangeUsername     command.ChangeUsernameHandler
	BanUser            command.BanUserHandler
	UnbanUser          command.UnbanUserHandler
}

type QueryHandler struct {
	GetUserById    query.GetUserByIdHandler
	GetUsers       query.GetUsersHandler
	GetUserByEmail query.GetUserByEmailHandler
}
