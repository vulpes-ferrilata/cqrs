package eventbus

import (
	"context"
	"reflect"
	"sync"

	"golang.org/x/sync/errgroup"
)

type key struct {
	eventPkgPath string
	eventName    string
}

type EventMiddlewareFunc func(eventHandlerFunc EventHandlerFunc[any]) EventHandlerFunc[any]

type EventBus interface {
	Use(middlewareFunc EventMiddlewareFunc) error
	Register(handlerFunc interface{}) error
	Publish(ctx context.Context, events ...interface{}) error
}

func NewEventBus() EventBus {
	return &eventBus{
		handlerFuncs:    make(map[key][]EventHandlerFunc[any]),
		middlewareFuncs: make([]EventMiddlewareFunc, 0),
	}
}

type eventBus struct {
	mu              sync.RWMutex
	handlerFuncs    map[key][]EventHandlerFunc[any]
	middlewareFuncs []EventMiddlewareFunc
}

func (e *eventBus) validateMiddlewareFunc(middlewareFunc EventMiddlewareFunc) error {
	reflectedMiddlewareFunc := reflect.ValueOf(middlewareFunc)

	if reflectedMiddlewareFunc.IsNil() {
		return ErrMiddlewareFuncMustNotBeNil
	}

	return nil
}

func (e *eventBus) Use(middlewareFunc EventMiddlewareFunc) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.validateMiddlewareFunc(middlewareFunc); err != nil {
		return err
	}

	e.middlewareFuncs = append(e.middlewareFuncs, middlewareFunc)

	return nil
}

func (e *eventBus) getHandlerFuncs(k key) []EventHandlerFunc[any] {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.handlerFuncs[k]
}

func (e *eventBus) setHandlerFunc(k key, handlerFunc EventHandlerFunc[any]) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.handlerFuncs[k] = append(e.handlerFuncs[k], handlerFunc)
}

func (e *eventBus) validateHandlerFunc(handlerFunc interface{}) error {
	reflectedHandlerFunc := reflect.ValueOf(handlerFunc)

	if reflectedHandlerFunc.Kind() != reflect.Func || reflectedHandlerFunc.IsNil() {
		return ErrHandlerFuncMustBeNonNilFunction
	}

	if reflectedHandlerFunc.Type().NumIn() != 2 {
		return ErrHandlerFuncMustHaveExactTwoArguments
	}

	firstArgumentType := reflectedHandlerFunc.Type().In(0)
	secondArgumentType := reflectedHandlerFunc.Type().In(1)

	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if firstArgumentType != contextType {
		return ErrFirstArgumentOfHandlerMustBeContext
	}

	if secondArgumentType.Kind() != reflect.Pointer || secondArgumentType.Elem().Kind() != reflect.Struct {
		return ErrSecondArgumentOfHandlerMustBePointerOfStruct
	}

	if reflectedHandlerFunc.Type().NumOut() != 1 {
		return ErrHandlerFuncMustHaveExactOneResult
	}

	firstResultType := reflectedHandlerFunc.Type().Out(0)
	errType := reflect.TypeOf((*error)(nil)).Elem()
	if firstResultType != errType {
		return ErrResultMustBeError
	}

	return nil
}

func (e *eventBus) getKeyFromHandlerFunc(handlerFunc interface{}) key {
	reflectedHandlerFunc := reflect.TypeOf(handlerFunc)
	secondArgumentType := reflectedHandlerFunc.In(1).Elem()

	return key{
		eventPkgPath: secondArgumentType.PkgPath(),
		eventName:    secondArgumentType.Name(),
	}
}

func (e *eventBus) wrapEventHandlerFunc(handlerFunc interface{}) EventHandlerFunc[any] {
	return func(ctx context.Context, event interface{}) (err error) {
		reflectedHandlerFunc := reflect.ValueOf(handlerFunc)
		reflectedArguments := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(event),
		}

		reflectedOut := reflectedHandlerFunc.Call(reflectedArguments)
		if !reflectedOut[0].IsNil() {
			err = reflectedOut[0].Interface().(error)
		}

		return
	}
}

func (e *eventBus) Register(handlerFunc interface{}) error {
	if err := e.validateHandlerFunc(handlerFunc); err != nil {
		return err
	}

	k := e.getKeyFromHandlerFunc(handlerFunc)

	e.setHandlerFunc(k, e.wrapEventHandlerFunc(handlerFunc))

	return nil
}

func (e *eventBus) validateEvent(event interface{}) error {
	reflectedEvent := reflect.ValueOf(event)

	if reflectedEvent.Kind() != reflect.Pointer || reflectedEvent.Elem().Kind() != reflect.Struct || reflectedEvent.IsNil() {
		return ErrEventMustBeNonNilPointerOfStruct
	}

	return nil
}

func (e *eventBus) getKeyFromEvent(event interface{}) key {
	reflectedEvent := reflect.TypeOf(event).Elem()

	return key{
		eventPkgPath: reflectedEvent.PkgPath(),
		eventName:    reflectedEvent.Name(),
	}
}

func (e *eventBus) Publish(ctx context.Context, events ...interface{}) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, event := range events {
		event := event

		if err := e.validateEvent(event); err != nil {
			return err
		}

		k := e.getKeyFromEvent(event)

		handlerFuncs := e.getHandlerFuncs(k)

		for _, handlerFunc := range handlerFuncs {
			handlerFunc := handlerFunc

			for i := len(e.middlewareFuncs) - 1; i >= 0; i-- {
				handlerFunc = e.middlewareFuncs[i](handlerFunc)
			}

			group.Go(func() error {
				return handlerFunc(ctx, event)
			})
		}
	}

	return group.Wait()
}
