package querybus

import (
	"context"
	"reflect"
	"sync"
)

type key struct {
	queryPkgPath string
	queryName    string
}

type QueryMiddlewareFunc func(queryHandlerFunc QueryHandlerFunc[any, any]) QueryHandlerFunc[any, any]

type QueryBus interface {
	Use(middlewareFunc QueryMiddlewareFunc) error
	Register(handlerFunc interface{}) error
	Execute(ctx context.Context, query interface{}) (interface{}, error)
}

func NewQueryBus() QueryBus {
	return &queryBus{
		handlerFuncs:    make(map[key]QueryHandlerFunc[any, any]),
		middlewareFuncs: make([]QueryMiddlewareFunc, 0),
	}
}

type queryBus struct {
	mu              sync.RWMutex
	handlerFuncs    map[key]QueryHandlerFunc[any, any]
	middlewareFuncs []QueryMiddlewareFunc
}

func (q *queryBus) validateMiddlewareFunc(middlewareFunc QueryMiddlewareFunc) error {
	reflectedMiddlewareFunc := reflect.ValueOf(middlewareFunc)

	if reflectedMiddlewareFunc.IsNil() {
		return ErrMiddlewareFuncMustNotBeNil
	}

	return nil
}

func (q *queryBus) Use(middlewareFunc QueryMiddlewareFunc) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if err := q.validateMiddlewareFunc(middlewareFunc); err != nil {
		return err
	}

	q.middlewareFuncs = append(q.middlewareFuncs, middlewareFunc)

	return nil
}

func (q *queryBus) getHandlerFunc(k key) QueryHandlerFunc[any, any] {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.handlerFuncs[k]
}

func (q *queryBus) setHandlerFunc(k key, handlerFunc QueryHandlerFunc[any, any]) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.handlerFuncs[k] = handlerFunc
}

func (q *queryBus) validateHandlerFunc(handlerFunc interface{}) error {
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

	if reflectedHandlerFunc.Type().NumOut() != 2 {
		return ErrHandlerFuncMustHaveExactTwoResults
	}

	firstResultType := reflectedHandlerFunc.Type().Out(0)
	secondResultType := reflectedHandlerFunc.Type().Out(1)

	if firstResultType.Kind() != reflect.Pointer || firstResultType.Elem().Kind() != reflect.Struct {
		return ErrFirstResultOfHandlerMustBePointerOfStruct
	}

	errType := reflect.TypeOf((*error)(nil)).Elem()
	if secondResultType != errType {
		return ErrSecondResultOfHandlerMustBeError
	}

	return nil
}

func (q *queryBus) getKeyFromHandlerFunc(handlerFunc interface{}) key {
	reflectedHandlerFunc := reflect.TypeOf(handlerFunc)
	secondArgumentType := reflectedHandlerFunc.In(1).Elem()

	return key{
		queryPkgPath: secondArgumentType.PkgPath(),
		queryName:    secondArgumentType.Name(),
	}
}

func (q *queryBus) wrapQueryHandlerFunc(handlerFunc interface{}) QueryHandlerFunc[any, any] {
	return func(ctx context.Context, query interface{}) (result interface{}, err error) {
		reflectedHandlerFunc := reflect.ValueOf(handlerFunc)
		reflectedArguments := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(query),
		}

		reflectedOut := reflectedHandlerFunc.Call(reflectedArguments)
		if !reflectedOut[0].IsNil() {
			result = reflectedOut[0].Interface()
		}
		if !reflectedOut[1].IsNil() {
			err = reflectedOut[1].Interface().(error)
		}

		return result, err
	}
}

func (q *queryBus) Register(handlerFunc interface{}) error {
	if err := q.validateHandlerFunc(handlerFunc); err != nil {
		return err
	}

	k := q.getKeyFromHandlerFunc(handlerFunc)

	if handlerFunc := q.getHandlerFunc(k); handlerFunc != nil {
		return ErrQueryIsAlreadyRegistered
	}

	q.setHandlerFunc(k, q.wrapQueryHandlerFunc(handlerFunc))

	return nil
}

func (q *queryBus) validateQuery(query interface{}) error {
	reflectedQuery := reflect.ValueOf(query)

	if reflectedQuery.Kind() != reflect.Pointer || reflectedQuery.Elem().Kind() != reflect.Struct || reflectedQuery.IsNil() {
		return ErrQueryMustBeNonNilPointerOfStruct
	}

	return nil
}

func (q *queryBus) getKeyFromQuery(query interface{}) key {
	reflectedQuery := reflect.TypeOf(query).Elem()

	return key{
		queryPkgPath: reflectedQuery.PkgPath(),
		queryName:    reflectedQuery.Name(),
	}
}

func (q *queryBus) Execute(ctx context.Context, query interface{}) (interface{}, error) {
	if err := q.validateQuery(query); err != nil {
		return nil, err
	}

	k := q.getKeyFromQuery(query)

	handlerFunc := q.getHandlerFunc(k)
	if handlerFunc == nil {
		return nil, ErrQueryHasNotRegisteredYet
	}

	for i := len(q.middlewareFuncs) - 1; i >= 0; i-- {
		handlerFunc = q.middlewareFuncs[i](handlerFunc)
	}

	return handlerFunc(ctx, query)
}
