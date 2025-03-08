package shared

import "context"

type QueryHandler[T any, D any] interface {
	Handle(ctx context.Context, cmd T) (D, error)
}
