package domain

import "time"

type UserReadModel struct {
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	Role       UserRole       `json:"role"`
	Id         string         `json:"id"`
	Reputation UserReputation `json:"reputation"`
	CreatedAt  time.Time      `json:"createtAt"`
	UpdatedAt  time.Time      `json:"createdAt"`
}

type UserReputation struct {
	ReputationScore int      `json:"reputationScore"`
	Badges          []string `json:"badges"`
}
