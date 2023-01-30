package cqrs

import (
	"context"
	"fmt"
	"reflect"
)

type QueryMiddlewareFunc func(queryHandlerFunc QueryHandlerFunc[any, any]) QueryHandlerFunc[any, any]

type QueryBus struct {
	handlers    map[string]QueryHandlerFunc[any, any]
	middlewares []QueryMiddlewareFunc
}

func (q *QueryBus) Register(query interface{}, handler QueryHandlerFunc[any, any]) error {
	queryName := reflect.TypeOf(query).String()

	if q.handlers == nil {
		q.handlers = make(map[string]QueryHandlerFunc[any, any])
	}

	q.handlers[queryName] = handler

	return nil
}

func (q *QueryBus) Use(middlewares ...QueryMiddlewareFunc) {
	q.middlewares = append(q.middlewares, middlewares...)
}

func (q QueryBus) Execute(ctx context.Context, query interface{}) (interface{}, error) {
	queryName := reflect.TypeOf(query).String()

	handler, ok := q.handlers[queryName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrHandlerNotFound, queryName)
	}

	for i := len(q.middlewares) - 1; i >= 0; i-- {
		handler = q.middlewares[i](handler)
	}

	return handler(ctx, query)
}
