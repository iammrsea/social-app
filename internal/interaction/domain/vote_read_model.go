package domain

type VoteReadMoel struct {
	UserID string   `json:"user_id"`
	PostID string   `json:"post_id"`
	Type   VoteType `json:"type"`
}
