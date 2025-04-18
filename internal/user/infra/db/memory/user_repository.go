package memory

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/iammrsea/social-app/internal/user/domain"
)

// simulate users table for a typical sql db
type userModel struct {
	id         string
	email      string
	username   string
	role       string
	reputation userReputationModel
	createdAt  time.Time
}

// simulate user_reputations table for a typical sql db
type userReputationModel struct {
	reputationScore int
	badges          []string
}

type memoryRepository struct {
	users []*userModel
}

func NewUserRepository(ctx context.Context) *memoryRepository {
	return &memoryRepository{}
}

func (m *memoryRepository) GetUserById(ctx context.Context, userId string) (*domain.UserReadModel, error) {
	u, err := m.getUserModelById(userId)
	if err != nil {
		return nil, err
	}
	return &domain.UserReadModel{
		Username: u.username,
		Email:    u.email,
		Role:     u.role,
		Id:       u.id,
		Reputation: domain.UserReputation{
			ReputationScore: u.reputation.reputationScore,
			Badges:          u.reputation.badges,
		},
	}, nil
}

func (m *memoryRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserReadModel, error) {
	i := slices.IndexFunc(m.users, func(u *userModel) bool {
		return email == u.email
	})
	if i < 0 {
		return nil, fmt.Errorf("user with email %s does not exist", email)
	}
	u := m.users[i]
	return &domain.UserReadModel{
		Username: u.username,
		Email:    u.email,
		Role:     u.role,
		Id:       u.id,
		Reputation: domain.UserReputation{
			ReputationScore: u.reputation.reputationScore,
			Badges:          u.reputation.badges,
		},
	}, nil
}

func (m *memoryRepository) GetUsers(ctx context.Context, opts domain.GetUsersOptions) ([]*domain.UserReadModel, bool, error) {
	users := []*domain.UserReadModel{}

	for _, user := range m.users {
		users = append(users, &domain.UserReadModel{
			Username: user.username,
			Email:    user.email,
			Role:     user.role,
			Id:       user.id,
			Reputation: domain.UserReputation{
				ReputationScore: user.reputation.reputationScore,
				Badges:          user.reputation.badges,
			},
		})
	}
	hasNext := false // TODO: Calculate hasNext
	return users, hasNext, nil
}

func (m *memoryRepository) Register(ctx context.Context, user domain.User) error {
	userExists := slices.ContainsFunc(m.users, func(u *userModel) bool {
		return user.Email() == u.email || user.Username() == u.username
	})
	if userExists {
		return errors.New("user with email or username already exists")
	}
	m.users = append(m.users, m.toUserModel(&user))
	return nil
}

func (m *memoryRepository) MakeModerator(ctx context.Context, userId string, updateFn func(user *domain.User) (*domain.User, error)) error {
	u, err := m.getUserModelById(userId)

	if err != nil {
		return err
	}
	userDomain := m.toDomainUser(u)
	updatedUser, err := updateFn(userDomain)

	if err != nil {
		return err
	}

	u.role = string(updatedUser.Role())

	return nil
}
func (m *memoryRepository) AwardBadge(ctx context.Context, userId string, updateFn func(user *domain.User) (*domain.User, error)) error {
	u, err := m.getUserModelById(userId)

	if err != nil {
		return err
	}
	userDomain := m.toDomainUser(u)

	updatedUser, err := updateFn(userDomain)

	if err != nil {
		return err
	}
	u.reputation.badges = updatedUser.Badges()
	return nil
}
func (m *memoryRepository) RevokeAwardedBadge(ctx context.Context, userId string, updateFn func(user *domain.User) (*domain.User, error)) error {
	u, err := m.getUserModelById(userId)
	if err != nil {
		return err
	}
	userDomain := m.toDomainUser(u)
	updatedUser, err := updateFn(userDomain)
	if err != nil {
		return err
	}
	u.reputation.badges = updatedUser.Badges()

	return nil
}
func (m *memoryRepository) ChangeUsername(ctx context.Context, userId string, updateFn func(user *domain.User) (*domain.User, error)) error {
	u, err := m.getUserModelById(userId)
	if err != nil {
		return err
	}
	userDomain := m.toDomainUser(u)

	updatedUser, err := updateFn(userDomain)

	if err != nil {
		return err
	}
	u.username = updatedUser.Username()
	return nil
}

func (m *memoryRepository) toUserModel(user *domain.User) *userModel {
	return &userModel{
		id:        user.Id(),
		email:     user.Email(),
		username:  user.Username(),
		role:      string(user.Role()),
		createdAt: user.JoinedAt(),
		reputation: userReputationModel{
			reputationScore: user.ReputationScore(),
			badges:          user.Badges(),
		},
	}
}

func (m *memoryRepository) toDomainUser(userModel *userModel) *domain.User {
	user := domain.MustNewUser(
		userModel.id,
		userModel.email,
		userModel.username,
		domain.UserRole(userModel.role),
		domain.MustNewUserReputation(userModel.reputation.reputationScore, userModel.reputation.badges),
	)

	return &user
}

func (m *memoryRepository) getUserModelById(userId string) (*userModel, error) {
	i := slices.IndexFunc(m.users, func(u *userModel) bool {
		return userId == u.id
	})
	if i < 0 {
		return nil, fmt.Errorf("user with id %s does not exist", userId)
	}
	return m.users[i], nil
}
