package cqrs

import (
	"context"
	"fmt"
	"reflect"
)

type EventMiddlewareFunc func(EventHandlerFunc EventHandlerFunc[any]) EventHandlerFunc[any]

type EventBus struct {
	handlers    map[string][]EventHandlerFunc[any]
	middlewares []EventMiddlewareFunc
}

func (e *EventBus) Register(event interface{}, handler EventHandlerFunc[any]) error {
	eventName := reflect.TypeOf(event).String()

	if e.handlers == nil {
		e.handlers = make(map[string][]EventHandlerFunc[any])
	}

	e.handlers[eventName] = append(e.handlers[eventName], handler)

	return nil
}

func (e *EventBus) Use(middlewares ...EventMiddlewareFunc) {
	e.middlewares = append(e.middlewares, middlewares...)
}

func (e EventBus) Publish(ctx context.Context, event interface{}) error {
	eventName := reflect.TypeOf(event).String()

	handlers, ok := e.handlers[eventName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrHandlerNotFound, eventName)
	}

	for _, handler := range handlers {
		for i := len(e.middlewares) - 1; i >= 0; i-- {
			handler = e.middlewares[i](handler)
		}

		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
