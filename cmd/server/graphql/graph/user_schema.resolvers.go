package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.70

import (
	"context"
	"time"

	"github.com/iammrsea/social-app/cmd/server/graphql/graph/model"
	"github.com/iammrsea/social-app/internal/shared/pagination"
	"github.com/iammrsea/social-app/internal/user/app/command"
	"github.com/iammrsea/social-app/internal/user/app/query"
	"github.com/iammrsea/social-app/internal/user/domain"
)

// ChangeUsername is the resolver for the changeUsername field.
func (r *mutationResolver) ChangeUsername(ctx context.Context, input model.ChangeUsername) (*domain.UserReadModel, error) {
	err := r.Services.UserService.CommandHandler.ChangeUsername.Handle(ctx, command.ChangeUsername{
		Id:       input.ID,
		Username: input.Username,
	})
	if err != nil {
		return nil, err
	}
	return r.Services.UserService.QueryHandler.GetUserById.Handle(ctx, query.GetUserById{
		Id: input.ID,
	})
}

// MakeModerator is the resolver for the makeModerator field.
func (r *mutationResolver) MakeModerator(ctx context.Context, id string) (*domain.UserReadModel, error) {
	err := r.Services.UserService.CommandHandler.MakeModerator.Handle(ctx, command.MakeModerator{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return r.Services.UserService.QueryHandler.GetUserById.Handle(ctx, query.GetUserById{
		Id: id,
	})
}

// BanUser is the resolver for the banUser field.
func (r *mutationResolver) BanUser(ctx context.Context, id string) (*domain.UserReadModel, error) {
	if err := r.Services.UserService.CommandHandler.BanUser.Handle(ctx, command.BanUser{
		Id: id,
	}); err != nil {
		return nil, err
	}
	return r.Services.UserService.QueryHandler.GetUserById.Handle(ctx, query.GetUserById{
		Id: id,
	})
}

// RegisterUser is the resolver for the registerUser field.
func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterUser) (*domain.UserReadModel, error) {
	err := r.Services.UserService.CommandHandler.RegisterUser.Handle(ctx, command.RegisterUser{
		Email:    input.Email,
		Username: input.Username,
	})
	if err != nil {
		return nil, err
	}
	return r.Services.UserService.QueryHandler.GetUserByEmail.Handle(ctx, query.GetUserByEmail{
		Email: input.Email,
	})
}

// AwardBadge is the resolver for the awardBadge field.
func (r *mutationResolver) AwardBadge(ctx context.Context, input model.AwardBadge) (*domain.UserReadModel, error) {
	err := r.Services.UserService.CommandHandler.AwardBadge.Handle(ctx, command.AwardBadge{
		Id:    input.ID,
		Badge: input.Badge,
	})
	if err != nil {
		return nil, err
	}
	return r.Services.UserService.QueryHandler.GetUserById.Handle(ctx, query.GetUserById{
		Id: input.ID,
	})
}

// RevokeAwardedBadge is the resolver for the revokeAwardedBadge field.
func (r *mutationResolver) RevokeAwardedBadge(ctx context.Context, input model.AwardBadge) (*domain.UserReadModel, error) {
	err := r.Services.UserService.CommandHandler.RevokeAwardedBadge.Handle(ctx, command.RevokeAwardedBadge{
		Id:    input.ID,
		Badge: input.Badge,
	})
	if err != nil {
		return nil, err
	}
	return r.Services.UserService.QueryHandler.GetUserById.Handle(ctx, query.GetUserById{
		Id: input.ID,
	})
}

// GetUserByID is the resolver for the getUserById field.
func (r *queryResolver) GetUserByID(ctx context.Context, id string) (*domain.UserReadModel, error) {
	return r.Services.UserService.QueryHandler.GetUserById.Handle(ctx, query.GetUserById{
		Id: id,
	})
}

// GetUsers is the resolver for the getUsers field.
func (r *queryResolver) GetUsers(ctx context.Context, first *int32, after *string) (*model.UserConnection, error) {
	var limit int32 = 10
	if first != nil {
		limit = *first
	}
	var afterCursor string

	if after != nil {
		decoded, err := pagination.DecodeCursor(*after)
		if err == nil {
			afterCursor = decoded
		}
	}

	result, err := r.Services.UserService.GetUsers.Handle(ctx, query.GetUsers{After: afterCursor, First: limit})

	if err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return &model.UserConnection{
			Edges:    []*model.UserEdge{},
			PageInfo: &pagination.PageInfo{},
		}, nil
	}

	edges := make([]*model.UserEdge, len(result.Data))

	for i, user := range result.Data {
		cursor := user.CreatedAt.UTC().Format(time.RFC3339Nano)
		edges[i] = &model.UserEdge{
			Cursor: pagination.EncodeCursor(cursor),
			Node:   user,
		}
	}
	startCursor := edges[0].Cursor
	endCursor := edges[len(edges)-1].Cursor

	return &model.UserConnection{
		Edges: edges,
		PageInfo: &pagination.PageInfo{
			HasNextPage:     result.PaginationInfo.HasNext,
			HasPreviousPage: afterCursor != "",
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
	}, nil
}

// GetUserByEmail is the resolver for the getUserByEmail field.
func (r *queryResolver) GetUserByEmail(ctx context.Context, email string) (*domain.UserReadModel, error) {
	return r.Services.UserService.QueryHandler.GetUserByEmail.Handle(ctx, query.GetUserByEmail{
		Email: email,
	})
}

// ReputationScore is the resolver for the reputationScore field.
func (r *userReputationResolver) ReputationScore(ctx context.Context, obj *domain.UserReputation) (int32, error) {
	return int32(obj.ReputationScore), nil
}

// UserReputation returns UserReputationResolver implementation.
func (r *Resolver) UserReputation() UserReputationResolver { return &userReputationResolver{r} }

type userReputationResolver struct{ *Resolver }
