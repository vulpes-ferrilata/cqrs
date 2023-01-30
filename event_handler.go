package cqrs

import "context"

type EventHandlerFunc[Event any] func(ctx context.Context, event Event) error

func WrapEventHandlerFunc[Event any](handler EventHandlerFunc[Event]) EventHandlerFunc[any] {
	return func(ctx context.Context, event interface{}) error {
		return handler(ctx, event.(Event))
	}
}
