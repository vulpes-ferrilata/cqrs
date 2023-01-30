package cqrs

import "context"

type CommandHandlerFunc[Command any] func(ctx context.Context, command Command) error

func WrapCommandHandlerFunc[Command any](handler CommandHandlerFunc[Command]) CommandHandlerFunc[any] {
	return func(ctx context.Context, command interface{}) error {
		return handler(ctx, command.(Command))
	}
}
