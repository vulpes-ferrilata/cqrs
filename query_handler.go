package cqrs

import "context"

type QueryHandlerFunc[Query any, Result any] func(ctx context.Context, query Query) (Result, error)

func WrapQueryHandlerFunc[Query any, Result any](handler QueryHandlerFunc[Query, Result]) QueryHandlerFunc[any, any] {
	return func(ctx context.Context, query interface{}) (interface{}, error) {
		return handler(ctx, query.(Query))
	}
}

func ParseQueryHandlerFunc[Query any, Result any](handler QueryHandlerFunc[any, any]) QueryHandlerFunc[Query, Result] {
	return func(ctx context.Context, query Query) (Result, error) {
		result, err := handler(ctx, query)

		return result.(Result), err
	}
}
