package middlewares

import (
	"context"

	"github.com/vulpes-ferrilata/cqrs"
)

func NewEventDispatcherMiddleware(eventBus cqrs.EventBus) *EventDispatcherMiddleware {
	return &EventDispatcherMiddleware{
		eventBus: eventBus,
	}
}

type EventDispatcherMiddleware struct {
	eventBus cqrs.EventBus
}

func (e EventDispatcherMiddleware) CommandMiddleware() cqrs.CommandMiddlewareFunc {
	return func(handler cqrs.CommandHandlerFunc[any]) cqrs.CommandHandlerFunc[any] {
		return func(ctx context.Context, command any) error {
			if err := handler(ctx, command); err != nil {
				return err
			}

			eventProvider, ok := cqrs.GetEventProvider(ctx)
			if !ok {
				return cqrs.ErrEventProviderNotFound
			}

			if err := e.eventBus.Dispatch(ctx, eventProvider.GetEvents()); err != nil {
				return err
			}

			return nil
		}
	}
}
