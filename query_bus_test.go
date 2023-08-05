package cqrs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vulpes-ferrilata/cqrs"
)

type (
	Query  struct{}
	Result struct{}
)

func Test_queryBus_Register(t *testing.T) {
	t.Parallel()

	type args struct {
		handler interface{}
	}
	tests := []struct {
		name    string
		prepare func(queryBus cqrs.QueryBus) error
		args    args
		wantErr error
	}{
		{
			name: "handler is not function",
			prepare: func(queryBus cqrs.QueryBus) error {
				return nil
			},
			args: args{
				handler: "",
			},
			wantErr: cqrs.ErrHandlerMustBeNonNilFunction,
		},
		{
			name: "handler has no argument",
			prepare: func(queryBus cqrs.QueryBus) error {
				return nil
			},
			args: args{
				handler: func() {},
			},
			wantErr: cqrs.ErrHandlerMustHaveExactTwoArguments,
		},
		{
			name: "handler have no context argument",
			prepare: func(queryBus cqrs.QueryBus) error {
				return nil
			},
			args: args{
				handler: func(i int, query Query) {},
			},
			wantErr: cqrs.ErrFirstArgumentOfHandlerMustBeContext,
		},
		{
			name: "handler has slice query argument",
			prepare: func(queryBus cqrs.QueryBus) error {
				return nil
			},
			args: args{
				handler: func(ctx context.Context, query []Query) {},
			},
			wantErr: cqrs.ErrSecondArgumentOfHandlerMustBeStructOrPointerOfStruct,
		},
		{
			name: "handler has no return",
			prepare: func(queryBus cqrs.QueryBus) error {
				return nil
			},
			args: args{
				handler: func(ctx context.Context, query Query) {},
			},
			wantErr: cqrs.ErrHandlerMustHaveExactTwoResults,
		},
		{
			name: "handler return bool",
			prepare: func(queryBus cqrs.QueryBus) error {
				return nil
			},
			args: args{
				handler: func(ctx context.Context, query Query) (Result, bool) {
					return Result{}, true
				},
			},
			wantErr: cqrs.ErrSecondResultOfHandlerMustBeError,
		},
		{
			name: "handler already registered",
			prepare: func(queryBus cqrs.QueryBus) error {
				if err := queryBus.Register(func(ctx context.Context, query Query) (Result, error) {
					return Result{}, nil
				}); err != nil {
					return err
				}

				return nil
			},
			args: args{
				handler: func(ctx context.Context, query Query) (Result, error) {
					return Result{}, nil
				},
			},
			wantErr: cqrs.ErrQueryAlreadyRegistered,
		},
		{
			name: "success",
			prepare: func(queryBus cqrs.QueryBus) error {
				return nil
			},
			args: args{
				handler: func(ctx context.Context, query Query) (Result, error) {
					return Result{}, nil
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryBus := cqrs.NewQueryBus()

			if err := tt.prepare(queryBus); assert.NoError(t, err) {
				err := queryBus.Register(tt.args.handler)
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func Test_queryBus_Execute(t *testing.T) {
	t.Parallel()

	var (
		Err = errors.New("error")
	)

	type args struct {
		ctx   context.Context
		query Query
	}
	type wants struct {
		result interface{}
		err    error
	}

	tests := []struct {
		name    string
		args    args
		prepare func(queryBus cqrs.QueryBus) error
		wants   wants
	}{
		{
			name: "query has not registered yet",
			args: args{
				ctx:   context.Background(),
				query: Query{},
			},
			prepare: func(queryBus cqrs.QueryBus) error {
				return nil
			},
			wants: wants{
				result: nil,
				err:    cqrs.ErrQueryHasNotRegisteredYet,
			},
		},
		{
			name: "query handler return error",
			args: args{
				ctx:   context.Background(),
				query: Query{},
			},
			prepare: func(queryBus cqrs.QueryBus) error {
				if err := queryBus.Register(func(ctx context.Context, query Query) (Result, error) {
					return Result{}, Err
				}); err != nil {
					return err
				}

				return nil
			},
			wants: wants{
				result: nil,
				err:    Err,
			},
		},
		{
			name: "middleware func return error",
			args: args{
				ctx:   context.Background(),
				query: Query{},
			},
			prepare: func(queryBus cqrs.QueryBus) error {
				queryBus.Use(func(handler cqrs.QueryHandlerFunc[any, any]) cqrs.QueryHandlerFunc[any, any] {
					return func(ctx context.Context, query interface{}) (interface{}, error) {
						return nil, Err
					}
				})

				if err := queryBus.Register(func(ctx context.Context, query Query) (Result, error) {
					return Result{}, nil
				}); err != nil {
					return err
				}

				return nil
			},
			wants: wants{
				result: nil,
				err:    Err,
			},
		},
		{
			name: "query handler return nil",
			args: args{
				ctx:   context.Background(),
				query: Query{},
			},
			prepare: func(queryBus cqrs.QueryBus) error {
				if err := queryBus.Register(func(ctx context.Context, query Query) (Result, error) {
					return Result{}, nil
				}); err != nil {
					return err
				}

				return nil
			},
			wants: wants{
				result: Result{},
				err:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryBus := cqrs.NewQueryBus()

			if err := tt.prepare(queryBus); assert.NoError(t, err) {
				result, err := queryBus.Execute(tt.args.ctx, tt.args.query)
				assert.ErrorIs(t, err, tt.wants.err)
				assert.Equal(t, tt.wants.result, result)
			}
		})
	}
}
