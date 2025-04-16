package query

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type GetUsersRepository interface {
	GetUsers(ctx context.Context) ([]domain.UserReadModel, error)
}

type GetUsersCommand struct {
}

type GetUsersHandler = shared.QueryHandler[GetUsersCommand, []domain.UserReadModel]

type getUsersCommandHandler struct {
	queryRepo GetUsersRepository
}

func NewGetUsersCommandHandler(queryRepo GetUsersRepository) GetUsersHandler {
	if queryRepo == nil {
		panic("nil user repository")
	}
	return &getUsersCommandHandler{queryRepo: queryRepo}
}

func (g *getUsersCommandHandler) Handle(ctx context.Context, cmd GetUsersCommand) ([]domain.UserReadModel, error) {
	users, err := g.queryRepo.GetUsers(ctx)

	if err != nil {
		return []domain.UserReadModel{}, errors.Unwrap(err)
	}
	return users, nil
}
