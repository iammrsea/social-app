package mongodb

import (
	"time"

	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
)

// userDocument represents how a user is stored in MongoDB
type userDocument struct {
	ID        string         `bson:"_id"`
	Email     string         `bson:"email"`
	Username  string         `bson:"username"`
	Role      string         `bson:"role"`
	Reputaion userReputation `bson:"reputation"`
	CreatedAt time.Time      `bson:"createdAt"`
	UpdatedAt time.Time      `bson:"updatedAt"`
	BanStatus userBanStatus  `bson:"banStatus"`
}

type userReputation struct {
	ReputationScore int      `bson:"reputationScore"`
	Badges          []string `bson:"badges"`
}

type userBanStatus struct {
	IsBanned        bool      `bson:"isBanned"`
	BannedAt        time.Time `bson:"bannedAt"`
	BanStartDate    time.Time `bson:"banStartDate"`
	BanEndDate      time.Time `bson:"banEndDate"`
	ReasonForBan    string    `bson:"reasonForBan"`
	IsBanIndefinite bool      `bson:"isBanIndefinite"`
}

// fromDomain converts a domain User to userDocument
func fromDomain(user domain.User) userDocument {
	return userDocument{
		ID:       user.Id(),
		Email:    user.Email(),
		Username: user.Username(),
		Role:     user.Role().String(),
		Reputaion: userReputation{
			Badges:          user.Badges(),
			ReputationScore: user.ReputationScore(),
		},
		CreatedAt: user.JoinedAt(),
		UpdatedAt: user.UpdatedAt(),
		BanStatus: userBanStatus{
			IsBanned:        user.IsBanned(),
			BannedAt:        user.BannedAt(),
			BanStartDate:    user.BanStartDate(),
			BanEndDate:      user.BanEndDate(),
			IsBanIndefinite: user.IsBanIndefinite(),
		},
	}
}

// toDomain converts a userDocument to domain User
func (u userDocument) toDomain() domain.User {
	return domain.MustNewUser(
		u.ID,
		u.Email,
		u.Username,
		rbac.UserRole(u.Role),
		u.CreatedAt,
		u.UpdatedAt,
		domain.MustNewUserReputation(u.Reputaion.ReputationScore, u.Reputaion.Badges),
		domain.NewBan(u.BanStatus.IsBanned, u.BanStatus.ReasonForBan, u.BanStatus.IsBanIndefinite, u.BanStatus.BanStartDate, u.BanStatus.BanEndDate, u.BanStatus.BannedAt),
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
			ReputationScore: doc.Reputaion.ReputationScore,
			Badges:          doc.Reputaion.Badges,
		},
		BanStatus: domain.BanStatus{
			IsBanned:        doc.BanStatus.IsBanned,
			BannedAt:        doc.BanStatus.BannedAt,
			BanStartDate:    doc.BanStatus.BanStartDate,
			BanEndDate:      doc.BanStatus.BanEndDate,
			ReasonForBan:    doc.BanStatus.ReasonForBan,
			IsBanIndefinite: doc.BanStatus.IsBanIndefinite,
		},
	}
}
