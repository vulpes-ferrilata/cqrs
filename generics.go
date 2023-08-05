package cqrs

import "context"

func RegisterCommandHandler[Command any](commandBus CommandBus, handler CommandHandlerFunc[Command]) error {
	if err := commandBus.Register(handler); err != nil {
		return err
	}

	return nil
}

func RegisterEventHandler[Event any](eventBus EventBus, handler EventHandlerFunc[Event]) error {
	if err := eventBus.Register(handler); err != nil {
		return err
	}

	return nil
}

func RegisterQueryHandler[Query any, Result any](queryBus QueryBus, handler QueryHandlerFunc[Query, Result]) error {
	if err := queryBus.Register(handler); err != nil {
		return err
	}

	return nil
}

func ExecuteQuery[Query any, Result any](queryBus QueryBus, ctx context.Context, query Query) (Result, error) {
	var emptyResult Result

	result, err := queryBus.Execute(ctx, query)
	if err != nil {
		return emptyResult, err
	}

	return result.(Result), nil
}
