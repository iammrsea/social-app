package adapters

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/iammrsea/social-app/internal/user/app/queries"
	"github.com/iammrsea/social-app/internal/user/domain"
)

// simulate users table for a typical sql db
type userModel struct {
	id         string
	email      string
	username   string
	role       string
	reputation userReputationModel
}

// simulate user_reputations table for a typical sql db
type userReputationModel struct {
	reputationScore int
	badges          []string
}

type memoryRepository struct {
	users []userModel
}

func NewMemoryRepository(ctx context.Context) *memoryRepository {
	return &memoryRepository{}
}

func (m *memoryRepository) GetUserById(ctx context.Context, userId string) (queries.User, error) {
	u, err := m.getUserModelById(userId)
	if err != nil {
		return queries.User{}, err
	}
	return queries.User{
		Username: u.email,
		Email:    u.username,
		Role:     u.role,
		Id:       u.id,
		Reputation: queries.UserReputation{
			ReputationScore: u.reputation.reputationScore,
			Badges:          u.reputation.badges,
		},
	}, nil
}

func (m *memoryRepository) GetUsers(ctx context.Context) ([]queries.User, error) {
	users := []queries.User{}

	for _, user := range m.users {
		users = append(users, queries.User{
			Username: user.username,
			Email:    user.email,
			Role:     user.role,
			Id:       user.id,
			Reputation: queries.UserReputation{
				ReputationScore: user.reputation.reputationScore,
				Badges:          user.reputation.badges,
			},
		})
	}
	return users, nil
}

func (m *memoryRepository) Register(ctx context.Context, user domain.User) error {
	userExists := slices.ContainsFunc(m.users, func(u userModel) bool {
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

	u.role = updatedUser.Role()

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

func (m *memoryRepository) toUserModel(user *domain.User) userModel {
	return userModel{
		id:       user.Id(),
		email:    user.Email(),
		username: user.Email(),
		role:     user.Role(),
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
	i := slices.IndexFunc(m.users, func(u userModel) bool {
		return userId == u.id
	})
	if i < 0 {
		return nil, fmt.Errorf("user with id %s does not exist", userId)
	}
	return &m.users[i], nil
}
