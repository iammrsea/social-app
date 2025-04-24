package domain

import (
	"time"

	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
)

type UserReadModel struct {
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	Role       rbac.UserRole  `json:"role"`
	Id         string         `json:"id"`
	Reputation UserReputation `json:"reputation"`
	CreatedAt  time.Time      `json:"createtAt"`
	UpdatedAt  time.Time      `json:"createdAt"`
	BanStatus  BanStatus      `json:"banStatus"`
}

type UserReputation struct {
	ReputationScore int      `json:"reputationScore"`
	Badges          []string `json:"badges"`
}

type BanStatus struct {
	IsBanned        bool      `json:"isBanned"`
	BannedAt        time.Time `json:"bannedAt"`
	BanStartDate    time.Time `json:"banStartDate"`
	BanEndDate      time.Time `json:"banEndDate"`
	ReasonForBan    string    `json:"reasonForBan"`
	IsBanIndefinite bool      `json:"isBanIndefinite"`
}
