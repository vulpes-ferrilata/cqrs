package cqrs

import "context"

type EventHandler[Event any] interface {
	Handle(ctx context.Context, event Event) error
}

type EventHandlerFunc[Event any] func(ctx context.Context, event Event) error

type EventMiddlewareFunc func(handlerFunc EventHandlerFunc[any]) EventHandlerFunc[any]
