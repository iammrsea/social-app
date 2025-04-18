package domain

type VoteType string

const (
	Upvote   VoteType = "upvote"
	Downvote VoteType = "downvote"
)

type Vote struct {
	UserID string
	PostID string
	Type   VoteType
}
