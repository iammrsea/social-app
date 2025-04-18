package query

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type UserByIdRepository interface {
	GetUserById(ctx context.Context, id string) (*domain.UserReadModel, error)
}

type GetUserByIdHandler = shared.QueryHandler[GetUserByIdCommand, *domain.UserReadModel]

type GetUserByIdCommand struct {
	Id string
}

type getUserByIdCommandHandler struct {
	queryRepo domain.UserReadModelRepository
}

func NewGetUserByIdCommandHandler(queryRepo domain.UserReadModelRepository) GetUserByIdHandler {
	if queryRepo == nil {
		panic("nil user repository")
	}
	return &getUserByIdCommandHandler{queryRepo: queryRepo}
}

func (g *getUserByIdCommandHandler) Handle(ctx context.Context, cmd GetUserByIdCommand) (*domain.UserReadModel, error) {
	authUser := auth.GetUserFromCtx(ctx)
	if !authUser.IsAuthenticated() {
		return nil, errors.New("unauthorized")
	}
	user, err := g.queryRepo.GetUserById(ctx, cmd.Id)

	if err != nil {
		return nil, errors.Unwrap(err)
	}
	return user, nil
}
