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
	command.RegisterUserHandler
	command.RevokeAwardedBagdeHandler
	command.AwardBadgeHandler
	command.MakeModeratorHandler
	command.ChangeUsernameHandler
}

type QueryHandler struct {
	query.GetUserHandler
	query.GetUsersHandler
}
