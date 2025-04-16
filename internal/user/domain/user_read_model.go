package domain

type UserReadModel struct {
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	Role       string         `json:"role"`
	Id         string         `json:"id"`
	Reputation UserReputation `json:"reputation,omitempty"`
}

type UserReputation struct {
	ReputationScore int      `json:"reputation_score"`
	Badges          []string `json:"badges"`
}
