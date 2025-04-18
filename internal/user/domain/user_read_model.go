package domain

import "time"

type UserReadModel struct {
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	Role       UserRole       `json:"role"`
	Id         string         `json:"id"`
	Reputation UserReputation `json:"reputation"`
	CreatedAt  time.Time      `json:"created_at"`
}

type UserReputation struct {
	ReputationScore int      `json:"reputation_score"`
	Badges          []string `json:"badges"`
}
