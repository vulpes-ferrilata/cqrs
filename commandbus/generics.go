package commandbus

import "context"

func RegisterCommandHandlerWithCommandBus[Command any](commandBus CommandBus, commandHandler CommandHandler[Command]) error {
	return commandBus.Register(commandHandler.Handle)
}

func RegisterCommandHandlerFuncWithCommandBus[Command any](commandBus CommandBus, commandHandlerFunc CommandHandlerFunc[*Command]) error {
	return commandBus.Register(commandHandlerFunc)
}

func ExecuteCommandWithCommandBus[Command any](commandBus CommandBus, ctx context.Context, command *Command) error {
	return commandBus.Execute(ctx, command)
}
