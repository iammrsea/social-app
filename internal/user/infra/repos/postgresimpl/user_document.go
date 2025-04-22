package postgresimpl

import (
	"errors"
	"time"

	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/jackc/pgx/v5"
)

// userDocument represents how a user is stored in Postgres
type userDocument struct {
	ID              string    `db:"id"`
	Email           string    `db:"email"`
	Username        string    `db:"username"`
	Role            string    `db:"role"`
	ReputationScore int       `db:"reputation_score"`
	Badges          []string  `db:"badges"`
	IsBanned        bool      `db:"is_banned"`
	BannedAt        time.Time `db:"banned_at"`
	BanStartDate    time.Time `db:"ban_start_date"`
	BanEndDate      time.Time `db:"ban_end_date"`
	ReasonForBan    string    `db:"reason_for_ban"`
	IsBanIndefinite bool      `db:"is_ban_indefinite"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// Helper function to determine the comparison operator based on sort direction
func getComparisonOperator(sortDirection string) string {
	if sortDirection == "ASC" {
		return ">"
	}
	return "<"
}

// toDomain converts a userDocument to a domain.User
func (u *userDocument) toDomain() domain.User {
	return domain.MustNewUser(
		u.ID,
		u.Email,
		u.Username,
		rbac.UserRole(u.Role),
		u.CreatedAt,
		u.UpdatedAt,
		domain.MustNewUserReputation(u.ReputationScore, u.Badges),
		domain.NewBan(
			u.IsBanned,
			u.ReasonForBan,
			u.IsBanIndefinite,
			u.BanStartDate,
			u.BanEndDate,
			u.BannedAt,
		),
	)
}

// documentToReadModel converts userDocument to UserReadModel
func documentToReadModel(doc userDocument) *domain.UserReadModel {
	return &domain.UserReadModel{
		Username:  doc.Username,
		Email:     doc.Email,
		Id:        doc.ID,
		Role:      rbac.UserRole(doc.Role),
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		Reputation: domain.UserReputation{
			ReputationScore: doc.ReputationScore,
			Badges:          doc.Badges,
		},
		BanStatus: domain.BanStatus{
			IsBanned:        doc.IsBanned,
			BannedAt:        doc.BannedAt,
			BanStartDate:    doc.BanStartDate,
			BanEndDate:      doc.BanEndDate,
			ReasonForBan:    doc.ReasonForBan,
			IsBanIndefinite: doc.IsBanIndefinite,
		},
	}
}

func scanUserRow(row pgx.Row, doc *userDocument) error {
	err := row.Scan(
		&doc.ID,
		&doc.Username,
		&doc.Email,
		&doc.Role,
		&doc.ReputationScore,
		&doc.Badges,
		&doc.IsBanned,
		&doc.BannedAt,
		&doc.BanStartDate,
		&doc.BanEndDate,
		&doc.ReasonForBan,
		&doc.IsBanIndefinite,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}
