package query

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type UserByIdRepository interface {
	GetUserById(ctx context.Context, id string) (domain.UserReadModel, error)
}

type GetUserHandler = shared.QueryHandler[GetUserByIdCommand, domain.UserReadModel]

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

func (g *getUserByIdCommandHandler) Handle(ctx context.Context, cmd GetUserByIdCommand) (domain.UserReadModel, error) {
	user, err := g.queryRepo.GetUserById(ctx, cmd.Id)

	if err != nil {
		return domain.UserReadModel{}, errors.Unwrap(err)
	}
	return user, nil
}
