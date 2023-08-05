package cqrs

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock_$GOPACKAGE
import (
	"context"
	"reflect"
	"sync"
)

type QueryBus interface {
	Use(middlewares ...QueryMiddlewareFunc)
	Register(handler interface{}) error
	Execute(ctx context.Context, query interface{}) (interface{}, error)
}

func NewQueryBus() QueryBus {
	return &queryBus{
		middlewares: make([]QueryMiddlewareFunc, 0),
		handlers:    make(map[reflect.Type]QueryHandlerFunc[any, any]),
	}
}

type queryBus struct {
	middlewares []QueryMiddlewareFunc
	handlers    map[reflect.Type]QueryHandlerFunc[any, any]
	mu          sync.RWMutex
}

func (c *queryBus) validate(handler interface{}) error {
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

	if handlerVal.Type().NumOut() != 2 {
		return ErrHandlerMustHaveExactTwoResults
	}

	secondResultType := handlerVal.Type().Out(1)
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if secondResultType != errorType {
		return ErrSecondResultOfHandlerMustBeError
	}

	_, ok := c.handlers[secondArgType]
	if ok {
		return ErrQueryAlreadyRegistered
	}

	return nil
}

func (c *queryBus) wrapHandler(handler interface{}) QueryHandlerFunc[any, any] {
	handlerVal := reflect.ValueOf(handler)

	return func(ctx context.Context, query interface{}) (interface{}, error) {
		args := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(query),
		}

		results := handlerVal.Call(args)
		if !results[1].IsNil() {
			return nil, results[1].Interface().(error)
		}

		return results[0].Interface(), nil
	}
}

func (c *queryBus) Use(middlewares ...QueryMiddlewareFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *queryBus) Register(handler interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.validate(handler); err != nil {
		return err
	}

	handlerType := reflect.TypeOf(handler)
	queryType := handlerType.In(1)

	c.handlers[queryType] = c.wrapHandler(handler)

	return nil
}

func (c *queryBus) Execute(ctx context.Context, query interface{}) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	queryType := reflect.TypeOf(query)

	handler, ok := c.handlers[queryType]
	if !ok {
		return nil, ErrQueryHasNotRegisteredYet
	}

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}

	result, err := handler(ctx, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}
