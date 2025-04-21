package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserReadModelRepository struct {
	db *pgxpool.Pool
}

func NewUserReadModelRepository(db *pgxpool.Pool) *UserReadModelRepository {
	return &UserReadModelRepository{db: db}
}

func (r *UserReadModelRepository) GetUsers(ctx context.Context, opts domain.GetUsersOptions) (users []*domain.UserReadModel, hasNext bool, err error) {
	// Validate sort direction
	sortDirection := "ASC" // Default to ascending
	if opts.SortDirection != "" {
		if opts.SortDirection != "ASC" && opts.SortDirection != "DESC" {
			return nil, false, errors.New("invalid sort direction, must be 'ASC' or 'DESC'")
		}
		sortDirection = opts.SortDirection
	}

	// Base query with dynamic ORDER BY
	query := fmt.Sprintf(`
        SELECT id, username, email, role, reputation_score, badges, is_banned, created_at, updated_at
        FROM users
        WHERE ($1::TIMESTAMP IS NULL OR created_at %s $1)
        ORDER BY created_at %s
        LIMIT $2
    `, getComparisonOperator(sortDirection), sortDirection)

	// Handle "after" timestamp
	var afterTimestamp *time.Time
	if opts.After != "" {
		parsedTime, err := time.Parse(time.RFC3339Nano, opts.After)
		if err != nil {
			return nil, false, errors.New("invalid after timestamp format")
		}
		afterTimestamp = &parsedTime
	}

	// Execute query
	rows, err := r.db.Query(ctx, query, afterTimestamp, opts.First+1) // Fetch one extra row to check for "hasNext"
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	// Parse results
	for rows.Next() {
		var user userDocument
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.ReputationScore,
			&user.Badges,
			&user.IsBanned,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, false, err
		}
		users = append(users, documentToReadModel(user))
	}

	// Determine if there is a next page
	hasNext = len(users) > int(opts.First)
	if hasNext {
		users = users[:opts.First] // Trim the extra row
	}

	return users, hasNext, nil
}

func (r *UserReadModelRepository) GetUserById(ctx context.Context, id string) (*domain.UserReadModel, error) {
	query := `
        SELECT id, username, email, role, reputation_score, badges, is_banned, banned_at,
            ban_start_date, ban_end_date, reason_for_ban, is_ban_indefinite, created_at, updated_at
        FROM users WHERE id = $1
    `
	row := r.db.QueryRow(ctx, query, id)

	var doc userDocument

	err := scanUserRow(row, &doc)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	user := documentToReadModel(doc)
	return user, nil
}

func (r *UserReadModelRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserReadModel, error) {
	query := `
        SELECT id, username, email, role, reputation_score, badges, is_banned, banned_at,
            ban_start_date, ban_end_date, reason_for_ban, is_ban_indefinite, created_at, updated_at
        FROM users WHERE email = $1
    `
	row := r.db.QueryRow(ctx, query, email)
	var doc userDocument
	err := scanUserRow(row, &doc)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	user := documentToReadModel(doc)
	return user, nil
}
