package query

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/guards"
	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type GetUserByIdHandler = shared.QueryHandler[GetUserById, *domain.UserReadModel]

type GetUserById struct {
	Id string
}

type getUserByIdHandler struct {
	queryRepo domain.UserReadModelRepository
	guard     guards.Guards
}

func NewGetUserByIdHandler(queryRepo domain.UserReadModelRepository, guard guards.Guards) GetUserByIdHandler {
	if queryRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &getUserByIdHandler{queryRepo: queryRepo, guard: guard}
}

func (g *getUserByIdHandler) Handle(ctx context.Context, cmd GetUserById) (*domain.UserReadModel, error) {
	authUser := auth.GetUserFromCtx(ctx)
	if err := g.guard.Authorize(authUser.Role, rbac.ViewUser); err != nil {
		return nil, err
	}
	return g.queryRepo.GetUserById(ctx, cmd.Id)
}
