package cqrs

import (
	"context"
	"fmt"
	"reflect"
)

type CommandMiddlewareFunc func(commandHandlerFunc CommandHandlerFunc[any]) CommandHandlerFunc[any]

type CommandBus struct {
	handlers    map[string]CommandHandlerFunc[any]
	middlewares []CommandMiddlewareFunc
}

func (c *CommandBus) Register(command interface{}, handler CommandHandlerFunc[any]) {
	commandName := reflect.TypeOf(command).String()

	if c.handlers == nil {
		c.handlers = make(map[string]CommandHandlerFunc[any])
	}

	c.handlers[commandName] = handler
}

func (c *CommandBus) Use(middlewares ...CommandMiddlewareFunc) {
	c.middlewares = append(c.middlewares, middlewares...)
}

func (c CommandBus) Execute(ctx context.Context, command interface{}) error {
	commandName := reflect.TypeOf(command).String()

	handler, ok := c.handlers[commandName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrHandlerNotFound, commandName)
	}

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}

	return handler(ctx, command)
}
