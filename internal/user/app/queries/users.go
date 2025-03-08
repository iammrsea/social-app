package queries

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
)

type GetUsersRepository interface {
	GetUsers(ctx context.Context) ([]User, error)
}

type GetUsersCommand struct {
}

type GetUsersHandler = shared.QueryHandler[GetUsersCommand, []User]

type getUsersCommandHandler struct {
	queryRepo GetUsersRepository
}

func NewGetUsersCommandHandler(queryRepo GetUsersRepository) GetUsersHandler {
	if queryRepo == nil {
		panic("nil user repository")
	}
	return &getUsersCommandHandler{queryRepo: queryRepo}
}

func (g *getUsersCommandHandler) Handle(ctx context.Context, cmd GetUsersCommand) ([]User, error) {
	users, err := g.queryRepo.GetUsers(ctx)

	if err != nil {
		return []User{}, errors.Unwrap(err)
	}
	return users, nil
}
