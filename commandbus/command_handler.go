package commandbus

import "context"

type CommandHandler[Command any] interface {
	Handle(ctx context.Context, command *Command) error
}

type CommandHandlerFunc[Command any] func(ctx context.Context, command Command) error
