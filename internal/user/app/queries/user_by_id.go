package queries

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
)

type UserByIdRepository interface {
	GetUserById(ctx context.Context, id string) (User, error)
}

type GetUserHandler = shared.QueryHandler[GetUserByIdCommand, User]

type GetUserByIdCommand struct {
	Id string
}

type getUserByIdCommandHandler struct {
	queryRepo UserByIdRepository
}

func NewGetUserByIdCommandHandler(queryRepo UserByIdRepository) GetUserHandler {
	if queryRepo == nil {
		panic("nil user repository")
	}
	return &getUserByIdCommandHandler{queryRepo: queryRepo}
}

func (g *getUserByIdCommandHandler) Handle(ctx context.Context, cmd GetUserByIdCommand) (User, error) {
	user, err := g.queryRepo.GetUserById(ctx, cmd.Id)

	if err != nil {
		return User{}, errors.Unwrap(err)
	}
	return user, nil
}
