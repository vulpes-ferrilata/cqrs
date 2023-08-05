package cqrs

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock_$GOPACKAGE
import (
	"context"
	"reflect"
	"sync"

	"golang.org/x/sync/errgroup"
)

type EventBus interface {
	Use(middlewares ...EventMiddlewareFunc)
	Register(handler interface{}) error
	Dispatch(ctx context.Context, events []interface{}) error
}

func NewEventBus() EventBus {
	return &eventBus{
		middlewares: make([]EventMiddlewareFunc, 0),
		handlers:    make(map[reflect.Type][]EventHandlerFunc[any]),
	}
}

type eventBus struct {
	middlewares []EventMiddlewareFunc
	handlers    map[reflect.Type][]EventHandlerFunc[any]
	mu          sync.RWMutex
}

func (c *eventBus) validate(handler interface{}) error {
	handlerVal := reflect.ValueOf(handler)

	if handlerVal.Kind() != reflect.Func || handlerVal.IsNil() {
		return ErrHandlerMustBeNonNilFunction
	}

	if handlerVal.Type().NumIn() != 2 {
		return ErrHandlerMustHaveExactTwoArguments
	}

	firstArgType := handlerVal.Type().In(0)
	secondArgType := handlerVal.Type().In(1)

	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if firstArgType != contextType {
		return ErrFirstArgumentOfHandlerMustBeContext
	}

	if (secondArgType.Kind() != reflect.Pointer || secondArgType.Elem().Kind() != reflect.Struct) && secondArgType.Kind() != reflect.Struct {
		return ErrSecondArgumentOfHandlerMustBeStructOrPointerOfStruct
	}

	if handlerVal.Type().NumOut() != 1 {
		return ErrHandlerMustHaveExactOneResult
	}

	firstResultType := handlerVal.Type().Out(0)
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if firstResultType != errorType {
		return ErrHandlerResultMustBeError
	}

	return nil
}

func (c *eventBus) wrapHandler(handler interface{}) EventHandlerFunc[any] {
	handlerVal := reflect.ValueOf(handler)

	return func(ctx context.Context, event interface{}) error {
		args := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(event),
		}

		results := handlerVal.Call(args)
		if !results[0].IsNil() {
			return results[0].Interface().(error)
		}

		return nil
	}
}

func (c *eventBus) Use(middlewares ...EventMiddlewareFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *eventBus) Register(handler interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.validate(handler); err != nil {
		return err
	}

	handlerType := reflect.TypeOf(handler)
	eventType := handlerType.In(1)

	c.handlers[eventType] = append(c.handlers[eventType], c.wrapHandler(handler))

	return nil
}

func (c *eventBus) Dispatch(ctx context.Context, events []interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	wg, ctx := errgroup.WithContext(ctx)
	wg.SetLimit(-1)

	for _, event := range events {
		event := event
		eventType := reflect.TypeOf(event)

		handlers, ok := c.handlers[eventType]
		if !ok {
			continue
		}

		for _, handler := range handlers {
			handler := handler

			for i := len(c.middlewares) - 1; i >= 0; i-- {
				handler = c.middlewares[i](handler)
			}

			wg.Go(func() error {
				if err := handler(ctx, event); err != nil {
					return err
				}

				return nil
			})
		}
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	return nil
}
