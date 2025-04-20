package query

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/rbac"
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
	guard     rbac.RequestGuard
}

func NewGetUserByEmailCommandHandler(queryRepo domain.UserReadModelRepository, guard rbac.RequestGuard) GetUserByEmailHandler {
	if queryRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &getUserByEmailCommandHandler{queryRepo: queryRepo, guard: guard}
}

func (g *getUserByEmailCommandHandler) Handle(ctx context.Context, cmd GetUserByEmailCommand) (*domain.UserReadModel, error) {
	authUser := auth.GetUserFromCtx(ctx)
	if err := g.guard.Authorize(authUser.Role, rbac.ViewUser); err != nil {
		return nil, err
	}
	user, err := g.queryRepo.GetUserByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
