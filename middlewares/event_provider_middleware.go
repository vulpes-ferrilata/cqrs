package middlewares

import (
	"context"

	"github.com/vulpes-ferrilata/cqrs"
)

func NewEventProviderMiddleware() *EventProviderMiddleware {
	return &EventProviderMiddleware{}
}

type EventProviderMiddleware struct{}

func (e EventProviderMiddleware) CommandMiddleware() cqrs.CommandMiddlewareFunc {
	return func(handler cqrs.CommandHandlerFunc[any]) cqrs.CommandHandlerFunc[any] {
		return func(ctx context.Context, command any) error {
			eventProvider := cqrs.NewEventProvider()
			ctx = cqrs.WithEventProvider(ctx, eventProvider)

			if err := handler(ctx, command); err != nil {
				return err
			}

			return nil
		}
	}
}
