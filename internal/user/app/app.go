package app

import (
	"github.com/iammrsea/social-app/internal/user/app/commands"
	"github.com/iammrsea/social-app/internal/user/app/queries"
)

type Application struct {
	CommandHandlers
	QueryHandlers
}

type CommandHandlers struct {
	commands.RegisterUserHandler
	commands.RevokeAwardedBagdeHandler
	commands.AwardBadgeHandler
	commands.MakeModeratorHandler
	commands.ChangeUsernameHandler
}

type QueryHandlers struct {
	queries.GetUserHandler
	queries.GetUsersHandler
}
