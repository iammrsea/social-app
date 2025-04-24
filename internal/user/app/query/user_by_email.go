package query

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/guards"
	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type UserByEmailRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.UserReadModel, error)
}

type GetUserByEmailHandler = shared.QueryHandler[GetUserByEmail, *domain.UserReadModel]

type GetUserByEmail struct {
	Email string
}

type getUserByEmailHandler struct {
	queryRepo domain.UserReadModelRepository
	guard     guards.Guards
}

func NewGetUserByEmailHandler(queryRepo domain.UserReadModelRepository, guard guards.Guards) GetUserByEmailHandler {
	if queryRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &getUserByEmailHandler{queryRepo: queryRepo, guard: guard}
}

func (g *getUserByEmailHandler) Handle(ctx context.Context, cmd GetUserByEmail) (*domain.UserReadModel, error) {
	authUser := auth.GetUserFromCtx(ctx)
	if err := g.guard.Authorize(authUser.Role, rbac.ViewUser); err != nil {
		return nil, err
	}
	return g.queryRepo.GetUserByEmail(ctx, cmd.Email)
}
