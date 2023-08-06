package cqrs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/vulpes-ferrilata/cqrs"
	mock_cqrs "github.com/vulpes-ferrilata/cqrs/mocks"
)

func TestRegisterCommandHandler(t *testing.T) {
	t.Parallel()

	var (
		Err = errors.New("error")
	)

	type mocks struct {
		commandBus *mock_cqrs.MockCommandBus
	}
	type args struct {
		handler cqrs.CommandHandlerFunc[Command]
	}
	tests := []struct {
		name    string
		prepare func(mocks mocks, args args)
		args    args
		wantErr error
	}{
		{
			name: "register handler return error",
			prepare: func(mocks mocks, args args) {
				mocks.commandBus.EXPECT().Register(gomock.AssignableToTypeOf(args.handler)).Return(Err)
			},
			args: args{
				handler: func(ctx context.Context, command Command) error {
					return nil
				},
			},
			wantErr: Err,
		},
		{
			name: "success",
			prepare: func(mocks mocks, args args) {
				mocks.commandBus.EXPECT().Register(gomock.AssignableToTypeOf(args.handler)).Return(nil)
			},
			args: args{
				handler: func(ctx context.Context, command Command) error {
					return nil
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mocks := mocks{
				commandBus: mock_cqrs.NewMockCommandBus(mockCtrl),
			}

			tt.prepare(mocks, tt.args)

			err := cqrs.RegisterCommandHandler(mocks.commandBus, tt.args.handler)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRegisterEventHandler(t *testing.T) {
	t.Parallel()

	var (
		Err = errors.New("error")
	)

	type mocks struct {
		eventBus *mock_cqrs.MockEventBus
	}
	type args struct {
		handler cqrs.EventHandlerFunc[Event]
	}
	tests := []struct {
		name    string
		prepare func(mocks mocks, args args)
		args    args
		wantErr error
	}{
		{
			name: "register handler return error",
			prepare: func(mocks mocks, args args) {
				mocks.eventBus.EXPECT().Register(gomock.AssignableToTypeOf(args.handler)).Return(Err)
			},
			args: args{
				handler: func(ctx context.Context, event Event) error {
					return nil
				},
			},
			wantErr: Err,
		},
		{
			name: "success",
			prepare: func(mocks mocks, args args) {
				mocks.eventBus.EXPECT().Register(gomock.AssignableToTypeOf(args.handler)).Return(nil)
			},
			args: args{
				handler: func(ctx context.Context, event Event) error {
					return nil
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mocks := mocks{
				eventBus: mock_cqrs.NewMockEventBus(mockCtrl),
			}

			tt.prepare(mocks, tt.args)

			err := cqrs.RegisterEventHandler(mocks.eventBus, tt.args.handler)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRegisterQueryHandler(t *testing.T) {
	t.Parallel()

	var (
		Err = errors.New("error")
	)

	type mocks struct {
		queryBus *mock_cqrs.MockQueryBus
	}
	type args struct {
		handler cqrs.QueryHandlerFunc[Query, Result]
	}
	tests := []struct {
		name    string
		prepare func(mocks mocks, args args)
		args    args
		wantErr error
	}{
		{
			name: "register handler return error",
			prepare: func(mocks mocks, args args) {
				mocks.queryBus.EXPECT().Register(gomock.AssignableToTypeOf(args.handler)).Return(Err)
			},
			args: args{
				handler: func(ctx context.Context, query Query) (Result, error) {
					return Result{}, nil
				},
			},
			wantErr: Err,
		},
		{
			name: "success",
			prepare: func(mocks mocks, args args) {
				mocks.queryBus.EXPECT().Register(gomock.AssignableToTypeOf(args.handler)).Return(nil)
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
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mocks := mocks{
				queryBus: mock_cqrs.NewMockQueryBus(mockCtrl),
			}

			tt.prepare(mocks, tt.args)

			err := cqrs.RegisterQueryHandler(mocks.queryBus, tt.args.handler)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestExecuteQuery(t *testing.T) {
	t.Parallel()

	var (
		Err = errors.New("error")
	)

	type mocks struct {
		queryBus *mock_cqrs.MockQueryBus
	}
	type args struct {
		ctx   context.Context
		query Query
	}
	type wants struct {
		result Result
		err    error
	}
	tests := []struct {
		name    string
		prepare func(mocks mocks, args args, wants wants)
		args    args
		wants   wants
	}{
		{
			name: "handler return error",
			prepare: func(mocks mocks, args args, wants wants) {
				mocks.queryBus.EXPECT().Execute(args.ctx, args.query).Return(wants.result, wants.err)
			},
			args: args{
				ctx:   context.Background(),
				query: Query{},
			},
			wants: wants{
				result: Result{},
				err:    Err,
			},
		},
		{
			name: "success",
			prepare: func(mocks mocks, args args, wants wants) {
				mocks.queryBus.EXPECT().Execute(args.ctx, args.query).Return(wants.result, wants.err)
			},
			args: args{
				ctx:   context.Background(),
				query: Query{},
			},
			wants: wants{
				result: Result{},
				err:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mocks := mocks{
				queryBus: mock_cqrs.NewMockQueryBus(mockCtrl),
			}

			tt.prepare(mocks, tt.args, tt.wants)

			got, err := cqrs.ExecuteQuery[Query, Result](mocks.queryBus, tt.args.ctx, tt.args.query)
			assert.ErrorIs(t, err, tt.wants.err)
			assert.Equal(t, tt.wants.result, got)
		})
	}
}
