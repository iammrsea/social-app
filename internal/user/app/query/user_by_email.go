package query

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type UserByEmailRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.UserReadModel, error)
}

type GetUserByEmailHandler = shared.QueryHandler[GetUserByEmailCommand, *domain.UserReadModel]

type GetUserByEmailCommand struct {
	Email string
}

type getUserByEmailCommandHandler struct {
	queryRepo domain.UserReadModelRepository
}

func NewGetUserByEmailCommandHandler(queryRepo domain.UserReadModelRepository) GetUserByEmailHandler {
	if queryRepo == nil {
		panic("nil user repository")
	}
	return &getUserByEmailCommandHandler{queryRepo: queryRepo}
}

func (g *getUserByEmailCommandHandler) Handle(ctx context.Context, cmd GetUserByEmailCommand) (*domain.UserReadModel, error) {
	authUser := auth.GetUserFromCtx(ctx)
	if !authUser.IsAuthenticated() {
		return nil, errors.New("unauthorized")
	}
	user, err := g.queryRepo.GetUserByEmail(ctx, cmd.Email)

	if err != nil {
		return nil, errors.Unwrap(err)
	}
	return user, nil
}
