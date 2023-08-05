package cqrs

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock_$GOPACKAGE
import (
	"context"
	"reflect"
	"sync"
)

type CommandBus interface {
	Use(middlewares ...CommandMiddlewareFunc)
	Register(handler interface{}) error
	Execute(ctx context.Context, command interface{}) error
}

func NewCommandBus() CommandBus {
	return &commandBus{
		middlewares: make([]CommandMiddlewareFunc, 0),
		handlers:    make(map[reflect.Type]CommandHandlerFunc[any]),
	}
}

type commandBus struct {
	middlewares []CommandMiddlewareFunc
	handlers    map[reflect.Type]CommandHandlerFunc[any]
	mu          sync.RWMutex
}

func (c *commandBus) validate(handler interface{}) error {
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

	_, ok := c.handlers[secondArgType]
	if ok {
		return ErrCommandAlreadyRegistered
	}

	return nil
}

func (c *commandBus) wrapHandler(handler interface{}) CommandHandlerFunc[any] {
	handlerVal := reflect.ValueOf(handler)

	return func(ctx context.Context, command interface{}) error {
		args := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(command),
		}

		results := handlerVal.Call(args)
		if !results[0].IsNil() {
			return results[0].Interface().(error)
		}

		return nil
	}
}

func (c *commandBus) Use(middlewares ...CommandMiddlewareFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *commandBus) Register(handler interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.validate(handler); err != nil {
		return err
	}

	handlerType := reflect.TypeOf(handler)
	commandType := handlerType.In(1)

	c.handlers[commandType] = c.wrapHandler(handler)

	return nil
}

func (c *commandBus) Execute(ctx context.Context, command interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	commandType := reflect.TypeOf(command)

	handler, ok := c.handlers[commandType]
	if !ok {
		return ErrCommandHasNotRegisteredYet
	}

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}

	if err := handler(ctx, command); err != nil {
		return err
	}

	return nil
}
