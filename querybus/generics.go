package querybus

import "context"

func RegisterQueryHandlerWithQueryBus[Query any, Result any](queryBus QueryBus, queryHandler QueryHandler[Query, Result]) error {
	return queryBus.Register(queryHandler.Handle)
}

func RegisterQueryHandlerFuncWithQueryBus[Query any, Result any](queryBus QueryBus, queryHandlerFunc QueryHandlerFunc[*Query, *Result]) error {
	return queryBus.Register(queryHandlerFunc)
}

func ExecuteQueryWithQueryBus[Query any, Result any](queryBus QueryBus, ctx context.Context, query *Query) (*Result, error) {
	result, err := queryBus.Execute(ctx, query)
	if result == nil {
		return nil, err
	}

	return result.(*Result), err
}
