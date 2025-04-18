package query

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/shared"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/pagination"
	"github.com/iammrsea/social-app/internal/user/domain"
)

type GetUsersCommand = domain.GetUsersOptions

type Result = pagination.PaginatedQueryResult[[]*domain.UserReadModel]

type GetUsersHandler = shared.QueryHandler[GetUsersCommand, *Result]

type getUsersCommandHandler struct {
	queryRepo domain.UserReadModelRepository
}

func NewGetUsersCommandHandler(queryRepo domain.UserReadModelRepository) GetUsersHandler {
	if queryRepo == nil {
		panic("nil user repository")
	}
	return &getUsersCommandHandler{queryRepo: queryRepo}
}

func (g *getUsersCommandHandler) Handle(ctx context.Context, cmd GetUsersCommand) (*Result, error) {
	authUser := auth.GetUserFromCtx(ctx)
	if !authUser.IsAuthenticated() {
		return nil, errors.New("unauthorized")
	}
	if authUser.Role != domain.Admin && authUser.Role != domain.Moderator {
		return nil, errors.New("unauthorized")
	}
	users, hasNext, err := g.queryRepo.GetUsers(ctx, cmd)

	if err != nil {
		return nil, errors.Unwrap(err)
	}
	result := &Result{
		Data: users,
		PaginationInfo: &pagination.PagenationInfo{
			HasNext: hasNext,
		},
	}
	return result, nil
}
