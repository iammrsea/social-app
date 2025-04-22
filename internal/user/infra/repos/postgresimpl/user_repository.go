package postgresimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(ctx context.Context, user domain.User) error {
	query := `
        INSERT INTO users (
            id, username, email, role, reputation_score, badges, is_banned, banned_at,
            ban_start_date, ban_end_date, reason_for_ban, is_ban_indefinite, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `
	_, err := r.db.Exec(ctx, query,
		user.Id(),
		user.Username(),
		user.Email(),
		user.Role().String(),
		user.ReputationScore(),
		user.Badges(),
		user.IsBanned(),
		user.BannedAt(),
		user.BanStartDate(),
		user.BanEndDate(),
		user.ReasonForBan(),
		user.IsBanIndefinite(),
		user.JoinedAt(),
		user.UpdatedAt(),
	)
	return err
}

func (r *UserRepository) MakeModerator(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.updateUser(ctx, userId, updateFn)
}

func (r *UserRepository) AwardBadge(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.updateUser(ctx, userId, updateFn)
}

func (r *UserRepository) RevokeAwardedBadge(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.updateUser(ctx, userId, updateFn)
}

func (r *UserRepository) ChangeUsername(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.updateUser(ctx, userId, updateFn)
}

func (r *UserRepository) BanUser(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.updateUser(ctx, userId, updateFn)
}

func (r *UserRepository) UnbanUser(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.updateUser(ctx, userId, updateFn)
}

func (r *UserRepository) GetUserBy(ctx context.Context, fieldName string, value any) (*domain.User, error) {
	query := fmt.Sprintf(`
        SELECT id, username, email, role, reputation_score, badges, is_banned, banned_at, ban_start_date, ban_end_date,
            reason_for_ban, is_ban_indefinite, updated_at, created_at
        FROM users
        WHERE %s = $1
    `, fieldName)

	var doc userDocument
	row := r.db.QueryRow(ctx, query, value)
	err := scanUserRow(row, &doc)
	if err != nil {
		return nil, err
	}
	user := doc.toDomain()
	return &user, nil
}

func (r *UserRepository) updateUser(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	user, err := r.GetUserBy(ctx, "id", userId)
	if err != nil {
		return err
	}
	if err := updateFn(user); err != nil {
		return err
	}
	query := `
        UPDATE users
        SET username = $1, email = $2, role = $3, reputation_score = $4, badges = $5,
            is_banned = $6, banned_at = $7, ban_start_date = $8, ban_end_date = $9,
            reason_for_ban = $10, is_ban_indefinite = $11, updated_at = $12
        WHERE id = $13
    `
	_, err = r.db.Exec(ctx, query,
		user.Username(),
		user.Email(),
		user.Role().String(),
		user.ReputationScore(),
		user.Badges(),
		user.IsBanned(),
		user.BannedAt(),
		user.BanStartDate(),
		user.BanEndDate(),
		user.ReasonForBan(),
		user.IsBanIndefinite(),
		time.Now(),
		user.Id(),
	)
	return err
}
