package cqrs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vulpes-ferrilata/cqrs"
)

type (
	Event struct{}
)

func Test_eventBus_Register(t *testing.T) {
	t.Parallel()

	type args struct {
		handler interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "handler is not function",
			args: args{
				handler: "",
			},
			wantErr: cqrs.ErrHandlerMustBeNonNilFunction,
		},
		{
			name: "handler has one argument",
			args: args{
				handler: func(ctx context.Context) {},
			},
			wantErr: cqrs.ErrHandlerMustHaveExactTwoArguments,
		},
		{
			name: "handler have no context argument",
			args: args{
				handler: func(i int, event Event) {},
			},
			wantErr: cqrs.ErrFirstArgumentOfHandlerMustBeContext,
		},
		{
			name: "handler has slice event argument",
			args: args{
				handler: func(ctx context.Context, evt []Event) {},
			},
			wantErr: cqrs.ErrSecondArgumentOfHandlerMustBeStructOrPointerOfStruct,
		},
		{
			name: "handler has no return",
			args: args{
				handler: func(ctx context.Context, event Event) {},
			},
			wantErr: cqrs.ErrHandlerMustHaveExactOneResult,
		},
		{
			name: "handler return bool",
			args: args{
				handler: func(ctx context.Context, event Event) bool {
					return true
				},
			},
			wantErr: cqrs.ErrHandlerResultMustBeError,
		},
		{
			name: "success",
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
			eventBus := cqrs.NewEventBus()

			err := eventBus.Register(tt.args.handler)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_eventBus_Dispatch(t *testing.T) {
	t.Parallel()

	var (
		Err = errors.New("error")
	)

	type args struct {
		ctx    context.Context
		events []interface{}
	}

	tests := []struct {
		name    string
		args    args
		prepare func(eventBus cqrs.EventBus) error
		wantErr error
	}{
		{
			name: "event has not registered yet",
			args: args{
				ctx: context.Background(),
				events: []interface{}{
					Event{},
				},
			},
			prepare: func(eventBus cqrs.EventBus) error {
				return nil
			},
			wantErr: nil,
		},
		{
			name: "event handler return error",
			args: args{
				ctx: context.Background(),
				events: []interface{}{
					Event{},
				},
			},
			prepare: func(eventBus cqrs.EventBus) error {
				if err := eventBus.Register(func(ctx context.Context, event Event) error {
					return Err
				}); err != nil {
					return err
				}

				return nil
			},
			wantErr: Err,
		},
		{
			name: "middleware func return error",
			args: args{
				ctx: context.Background(),
				events: []interface{}{
					Event{},
				},
			},
			prepare: func(eventBus cqrs.EventBus) error {
				eventBus.Use(func(handler cqrs.EventHandlerFunc[any]) cqrs.EventHandlerFunc[any] {
					return func(ctx context.Context, event any) error {
						return Err
					}
				})

				if err := eventBus.Register(func(ctx context.Context, event Event) error {
					return nil
				}); err != nil {
					return err
				}

				return nil
			},
			wantErr: Err,
		},
		{
			name: "event handler return nil",
			args: args{
				ctx: context.Background(),
				events: []interface{}{
					Event{},
				},
			},
			prepare: func(eventBus cqrs.EventBus) error {
				if err := eventBus.Register(func(ctx context.Context, event Event) error {
					return nil
				}); err != nil {
					return err
				}

				return nil
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventBus := cqrs.NewEventBus()

			if err := tt.prepare(eventBus); assert.NoError(t, err) {
				err = eventBus.Dispatch(tt.args.ctx, tt.args.events)
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
