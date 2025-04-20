package query

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/pagination"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type GetUsers = domain.GetUsersOptions

type Result = pagination.PaginatedQueryResult[[]*domain.UserReadModel]

type GetUsersHandler = shared.QueryHandler[GetUsers, *Result]

type getUsersHandler struct {
	queryRepo domain.UserReadModelRepository
	guard     rbac.RequestGuard
}

func NewGetUsersHandler(queryRepo domain.UserReadModelRepository, guard rbac.RequestGuard) GetUsersHandler {
	if queryRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &getUsersHandler{queryRepo: queryRepo, guard: guard}
}

func (g *getUsersHandler) Handle(ctx context.Context, cmd GetUsers) (*Result, error) {
	authUser := auth.GetUserFromCtx(ctx)
	if err := g.guard.Authorize(authUser.Role, rbac.ListUsers); err != nil {
		return nil, err
	}
	users, hasNext, err := g.queryRepo.GetUsers(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &Result{
		Data: users,
		PaginationInfo: &pagination.PagenationInfo{
			HasNext: hasNext,
		},
	}, nil
}
