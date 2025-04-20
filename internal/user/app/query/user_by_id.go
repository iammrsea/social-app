package query

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type GetUserByIdHandler = shared.QueryHandler[GetUserByIdCommand, *domain.UserReadModel]

type GetUserByIdCommand struct {
	Id string
}

type getUserByIdCommandHandler struct {
	queryRepo domain.UserReadModelRepository
	guard     rbac.RequestGuard
}

func NewGetUserByIdCommandHandler(queryRepo domain.UserReadModelRepository, guard rbac.RequestGuard) GetUserByIdHandler {
	if queryRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &getUserByIdCommandHandler{queryRepo: queryRepo, guard: guard}
}

func (g *getUserByIdCommandHandler) Handle(ctx context.Context, cmd GetUserByIdCommand) (*domain.UserReadModel, error) {
	authUser := auth.GetUserFromCtx(ctx)
	if err := g.guard.Authorize(authUser.Role, rbac.ViewUser); err != nil {
		return nil, err
	}
	user, err := g.queryRepo.GetUserById(ctx, cmd.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
