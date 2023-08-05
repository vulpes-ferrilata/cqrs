package cqrs

import "context"

type QueryHandler[Query any, Result any] interface {
	Handle(ctx context.Context, query Query) (Result, error)
}

type QueryHandlerFunc[Query any, Result any] func(ctx context.Context, query Query) (Result, error)

type QueryMiddlewareFunc func(handlerFunc QueryHandlerFunc[any, any]) QueryHandlerFunc[any, any]
