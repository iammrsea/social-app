package query

import (
	"context"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/pagination"
	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type GetUsersCommand = domain.GetUsersOptions

type Result = pagination.PaginatedQueryResult[[]*domain.UserReadModel]

type GetUsersHandler = shared.QueryHandler[GetUsersCommand, *Result]

type getUsersCommandHandler struct {
	queryRepo domain.UserReadModelRepository
	guard     rbac.RequestGuard
}

func NewGetUsersCommandHandler(queryRepo domain.UserReadModelRepository, guard rbac.RequestGuard) GetUsersHandler {
	if queryRepo == nil || guard == nil {
		panic("nil user repository or guard")
	}
	return &getUsersCommandHandler{queryRepo: queryRepo, guard: guard}
}

func (g *getUsersCommandHandler) Handle(ctx context.Context, cmd GetUsersCommand) (*Result, error) {
	authUser := auth.GetUserFromCtx(ctx)
	if err := g.guard.Authorize(authUser.Role, rbac.ListUsers); err != nil {
		return nil, err
	}
	users, hasNext, err := g.queryRepo.GetUsers(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := &Result{
		Data: users,
		PaginationInfo: &pagination.PagenationInfo{
			HasNext: hasNext,
		},
	}
	return result, nil
}
