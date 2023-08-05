package middlewares_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/vulpes-ferrilata/cqrs"
	"github.com/vulpes-ferrilata/cqrs/middlewares"
	mock_cqrs "github.com/vulpes-ferrilata/cqrs/mocks"
)

func TestEventDispatcherMiddleware_CommandMiddleware(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.Background()
		command = struct{}{}
		events  = []interface{}{
			struct{}{},
			struct{}{},
		}
		Err = errors.New("error")
	)

	type mocks struct {
		eventBus *mock_cqrs.MockEventBus
	}
	type args struct {
		handler cqrs.CommandHandlerFunc[any]
		ctx     context.Context
		command interface{}
	}
	tests := []struct {
		name    string
		prepare func(mocks mocks, args args)
		args    args
		wantErr error
	}{
		{
			name:    "handler return error",
			prepare: func(mocks mocks, args args) {},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return Err
				},
				ctx:     cqrs.WithEventProvider(ctx, cqrs.NewEventProvider()),
				command: command,
			},
			wantErr: Err,
		},
		{
			name:    "event provider not found",
			prepare: func(mocks mocks, args args) {},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return nil
				},
				ctx:     ctx,
				command: command,
			},
			wantErr: cqrs.ErrEventProviderNotFound,
		},
		{
			name: "dispatch event fail",
			prepare: func(mocks mocks, args args) {
				mocks.eventBus.EXPECT().Dispatch(args.ctx, events).Return(Err)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					eventProvider, ok := cqrs.GetEventProvider(ctx)
					if ok {
						eventProvider.CollectEvents(events...)
					}

					return nil
				},
				ctx:     cqrs.WithEventProvider(ctx, cqrs.NewEventProvider()),
				command: command,
			},
			wantErr: Err,
		},
		{
			name: "success",
			prepare: func(mocks mocks, args args) {
				mocks.eventBus.EXPECT().Dispatch(args.ctx, events).Return(nil)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					eventProvider, ok := cqrs.GetEventProvider(ctx)
					if ok {
						eventProvider.CollectEvents(events...)
					}

					return nil
				},
				ctx:     cqrs.WithEventProvider(ctx, cqrs.NewEventProvider()),
				command: command,
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

			eventDispatcherMiddleware := middlewares.NewEventDispatcherMiddleware(mocks.eventBus)
			commandMiddleware := eventDispatcherMiddleware.CommandMiddleware()
			handler := commandMiddleware(tt.args.handler)
			err := handler(tt.args.ctx, tt.args.command)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
