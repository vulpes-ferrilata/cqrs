package commandbus

import (
	"context"
	"reflect"
	"sync"
)

type key struct {
	commandPkgPath string
	commandName    string
}

type CommandMiddlewareFunc func(commandHandlerFunc CommandHandlerFunc[any]) CommandHandlerFunc[any]

type CommandBus interface {
	Use(middlewareFunc CommandMiddlewareFunc) error
	Register(handlerFunc interface{}) error
	Execute(ctx context.Context, command interface{}) error
}

func NewCommandBus() CommandBus {
	return &commandBus{
		handlerFuncs:    make(map[key]CommandHandlerFunc[any]),
		middlewareFuncs: make([]CommandMiddlewareFunc, 0),
	}
}

type commandBus struct {
	mu              sync.RWMutex
	handlerFuncs    map[key]CommandHandlerFunc[any]
	middlewareFuncs []CommandMiddlewareFunc
}

func (c *commandBus) validateMiddlewareFunc(middlewareFunc CommandMiddlewareFunc) error {
	reflectedMiddlewareFunc := reflect.ValueOf(middlewareFunc)

	if reflectedMiddlewareFunc.IsNil() {
		return ErrMiddlewareFuncMustNotBeNil
	}

	return nil
}

func (c *commandBus) Use(middlewareFunc CommandMiddlewareFunc) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.validateMiddlewareFunc(middlewareFunc); err != nil {
		return err
	}

	c.middlewareFuncs = append(c.middlewareFuncs, middlewareFunc)

	return nil
}

func (c *commandBus) getHandlerFunc(k key) CommandHandlerFunc[any] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.handlerFuncs[k]
}

func (c *commandBus) setHandlerFunc(k key, handlerFunc CommandHandlerFunc[any]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlerFuncs[k] = handlerFunc
}

func (c *commandBus) validateHandlerFunc(handlerFunc interface{}) error {
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

func (c *commandBus) getKeyFromHandlerFunc(handlerFunc interface{}) key {
	reflectedHandlerFunc := reflect.TypeOf(handlerFunc)
	secondArgumentType := reflectedHandlerFunc.In(1).Elem()

	return key{
		commandPkgPath: secondArgumentType.PkgPath(),
		commandName:    secondArgumentType.Name(),
	}
}

func (c *commandBus) wrapCommandHandlerFunc(handlerFunc interface{}) CommandHandlerFunc[any] {
	return func(ctx context.Context, command interface{}) (err error) {
		reflectedHandlerFunc := reflect.ValueOf(handlerFunc)
		reflectedArguments := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(command),
		}

		reflectedOut := reflectedHandlerFunc.Call(reflectedArguments)
		if !reflectedOut[0].IsNil() {
			err = reflectedOut[0].Interface().(error)
		}

		return
	}
}

func (c *commandBus) Register(handlerFunc interface{}) error {
	if err := c.validateHandlerFunc(handlerFunc); err != nil {
		return err
	}

	k := c.getKeyFromHandlerFunc(handlerFunc)

	if handlerFunc := c.getHandlerFunc(k); handlerFunc != nil {
		return ErrCommandIsAlreadyRegistered
	}

	c.setHandlerFunc(k, c.wrapCommandHandlerFunc(handlerFunc))

	return nil
}

func (c *commandBus) validateCommand(command interface{}) error {
	reflectedCommand := reflect.ValueOf(command)

	if reflectedCommand.Kind() != reflect.Pointer || reflectedCommand.Elem().Kind() != reflect.Struct || reflectedCommand.IsNil() {
		return ErrCommandMustBeNonNilPointerOfStruct
	}

	return nil
}

func (c *commandBus) getKeyFromCommand(command interface{}) key {
	reflectedCommand := reflect.TypeOf(command).Elem()

	return key{
		commandPkgPath: reflectedCommand.PkgPath(),
		commandName:    reflectedCommand.Name(),
	}
}

func (c *commandBus) Execute(ctx context.Context, command interface{}) error {
	if err := c.validateCommand(command); err != nil {
		return err
	}

	k := c.getKeyFromCommand(command)

	handlerFunc := c.getHandlerFunc(k)
	if handlerFunc == nil {
		return ErrCommandHasNotRegisteredYet
	}

	for i := len(c.middlewareFuncs) - 1; i >= 0; i-- {
		handlerFunc = c.middlewareFuncs[i](handlerFunc)
	}

	return handlerFunc(ctx, command)
}
