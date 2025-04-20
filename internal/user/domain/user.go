package domain

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/iammrsea/social-app/internal/shared/rbac"
)

type User struct {
	id         string
	email      string
	username   string
	reputation *userReputation
	role       rbac.UserRole
	banStatus  *banning
	joinedAt   time.Time
	updatedAt  time.Time
}

type userReputation struct {
	reputationScore int
	badges          []string
}

var (
	ErrUserIdRequired        = errors.New("user id cannot be empty")
	ErrUserEmailRequired     = errors.New("user email cannot be empty")
	ErrUsernameRequired      = errors.New("username cannot be empty")
	ErrUserRoleRequired      = errors.New("user role cannot be empty")
	ErrBadgeRequired         = errors.New("badge cannot be empty")
	ErrInvalidIncrementValue = errors.New("you cannot increment user reputation by a value less than one")
	ErrInvalidDecrementValue = errors.New("you cannot decrement user reputation by a value less than one")
	ErrInvalidRepScore       = errors.New("invalid reputation score")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUserNotFound          = errors.New("user not found")
)

func NewUser(id, email, username string, role rbac.UserRole, joinedAt time.Time, updatedAt time.Time, reputation *userReputation) (User, error) {
	user := User{}
	if strings.TrimSpace(id) == "" {
		return user, ErrUserIdRequired
	}
	if strings.TrimSpace(email) == "" {
		return user, ErrUserEmailRequired
	}
	if strings.TrimSpace(username) == "" {
		return user, ErrUsernameRequired
	}
	if strings.TrimSpace(string(role)) == "" {
		return user, ErrUserRoleRequired
	}
	if !isValidRole(role) {
		return user, fmt.Errorf("invalid user role. Valid user roles are %s, %s and %s", rbac.Admin, rbac.Moderator, rbac.Regular)
	}

	if reputation == nil {
		reputation = &userReputation{
			reputationScore: 0,
			badges:          []string{},
		}
	}

	return User{id: id, email: email, joinedAt: joinedAt, updatedAt: updatedAt, username: username, role: role, reputation: reputation, banStatus: &banning{isBanned: false}}, nil
}

func MustNewUser(id, email, username string, role rbac.UserRole, joinedAt time.Time, updatedAt time.Time, reputation *userReputation) User {
	user, err := NewUser(id, email, username, role, joinedAt, updatedAt, reputation)
	if err != nil {
		panic(err.Error())
	}
	return user
}

func NewUserReputation(score int, badges []string) (*userReputation, error) {
	if score < 0 {
		return &userReputation{}, ErrInvalidRepScore
	}
	return &userReputation{reputationScore: score, badges: badges}, nil
}

func MustNewUserReputation(score int, badges []string) *userReputation {
	rep, err := NewUserReputation(score, badges)
	if err != nil {
		panic(err.Error())
	}
	return rep
}

func (u *User) ChangeUsername(newUsername string) error {
	if strings.TrimSpace(newUsername) == "" {
		return ErrUsernameRequired
	}
	u.username = newUsername
	u.updatedAt = time.Now()
	return nil
}

func (u *User) AwardBadge(badge string) error {
	if strings.TrimSpace(badge) == "" {
		return ErrBadgeRequired
	}
	u.reputation.badges = append(u.reputation.badges, badge)
	u.updatedAt = time.Now()
	return nil
}

func (u *User) RevokeAwardedBadge(badge string) error {
	if strings.TrimSpace(badge) == "" {
		return ErrBadgeRequired
	}
	badges := u.reputation.badges

	awardedBadge := slices.Contains(badges, badge)
	if !awardedBadge {
		return fmt.Errorf("the badge %s you want to revoke hasn't been awarded to the user %s previously", badge, u.username)
	}

	u.reputation.badges = slices.DeleteFunc(badges, func(awardedBadge string) bool {
		return awardedBadge == badge
	})
	u.updatedAt = time.Now()
	return nil
}

func (u *User) IncrementReputationScoreBy(v int) error {
	if v < 1 {
		return ErrInvalidIncrementValue
	}
	u.reputation.reputationScore = u.reputation.reputationScore + v
	u.updatedAt = time.Now()
	return nil
}

func (u *User) DecrementReputationScoreBy(v int) error {
	if v < 1 {
		return ErrInvalidDecrementValue
	}
	u.reputation.reputationScore = u.reputation.reputationScore - v
	u.updatedAt = time.Now()
	return nil
}

func (u *User) Id() string {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Username() string {
	return u.username
}

func (u *User) ReputationScore() int {
	return u.reputation.reputationScore
}

func (u *User) Badges() []string {
	return u.reputation.badges
}

func (u *User) JoinedAt() time.Time {
	return u.joinedAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
