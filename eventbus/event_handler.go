package eventbus

import "context"

type EventHandler[Event any] interface {
	Handle(ctx context.Context, event *Event) error
}

type EventHandlerFunc[Event any] func(ctx context.Context, event Event) error
