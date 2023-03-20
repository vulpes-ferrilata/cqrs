package querybus

import "context"

type QueryHandler[Query any, Result any] interface {
	Handle(ctx context.Context, query *Query) (*Result, error)
}

type QueryHandlerFunc[Query any, Result any] func(ctx context.Context, query Query) (Result, error)
