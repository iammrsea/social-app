package domain

import (
	"context"
)

type UserRepository interface {
	Register(ctx context.Context, user User) error
	MakeModerator(ctx context.Context, userId string, updateFn func(user *User) error) error
	AwardBadge(ctx context.Context, userId string, updateFn func(user *User) error) error
	RevokeAwardedBadge(ctx context.Context, userId string, updateFn func(user *User) error) error
	ChangeUsername(ctx context.Context, userId string, updateFn func(user *User) error) error
	UnbanUser(ctx context.Context, userId string, updateFn func(user *User) error) error
	BanUser(ctx context.Context, userId string, updateFn func(user *User) error) error
	GetUserBy(ctx context.Context, fieldName string, value any) (*User, error)
	UserExists(ctx context.Context, email string, username string) (bool, error)
}
