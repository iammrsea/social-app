package domain

import "context"

type GetUsersOptions struct {
	First         int32
	After         string
	SortDirection string // "ASC" or "DESC"
}

type UserReadModelRepository interface {
	GetUsers(ctx context.Context, opts GetUsersOptions) (users []*UserReadModel, hasNext bool, err error)
	GetUserById(ctx context.Context, id string) (*UserReadModel, error)
	GetUserByEmail(ctx context.Context, email string) (*UserReadModel, error)
}
