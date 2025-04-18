package pagination

import "encoding/base64"

type PageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
}

type PagenationInfo struct {
	HasNext bool
}
type PaginatedQueryResult[T any] struct {
	Data           T
	PaginationInfo *PagenationInfo
}

type Edge[T any] struct {
	Cursor string `json:"cursor"`
	Node   *T
}

type Connection[T any] struct {
	Edges    []*Edge[T] `json:"edges"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

func EncodeCursor(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(id))
}

func DecodeCursor(cursor string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)

	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
