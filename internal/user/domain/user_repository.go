package domain

import (
	"context"
)

type UserRepository interface {
	Register(ctx context.Context, user User) error
	MakeModerator(ctx context.Context, userId string, updateFn func(user *User) (*User, error)) error
	AwardBadge(ctx context.Context, userId string, updateFn func(user *User) (*User, error)) error
	RevokeAwardedBadge(ctx context.Context, userId string, updateFn func(user *User) (*User, error)) error
	ChangeUsername(ctx context.Context, userId string, updateFn func(user *User) (*User, error)) error
}
