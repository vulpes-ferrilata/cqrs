package cqrs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vulpes-ferrilata/cqrs"
)

type (
	Command struct{}
)

func Test_commandBus_Register(t *testing.T) {
	t.Parallel()

	type args struct {
		handler interface{}
	}
	tests := []struct {
		name    string
		prepare func(commandBus cqrs.CommandBus) error
		args    args
		wantErr error
	}{
		{
			name: "handler is not function",
			prepare: func(commandBus cqrs.CommandBus) error {
				return nil
			},
			args: args{
				handler: "",
			},
			wantErr: cqrs.ErrHandlerMustBeNonNilFunction,
		},
		{
			name: "handler has no argument",
			prepare: func(commandBus cqrs.CommandBus) error {
				return nil
			},
			args: args{
				handler: func() {},
			},
			wantErr: cqrs.ErrHandlerMustHaveExactTwoArguments,
		},
		{
			name: "handler have no context argument",
			prepare: func(commandBus cqrs.CommandBus) error {
				return nil
			},
			args: args{
				handler: func(i int, command Command) {},
			},
			wantErr: cqrs.ErrFirstArgumentOfHandlerMustBeContext,
		},
		{
			name: "handler has slice command argument",
			prepare: func(commandBus cqrs.CommandBus) error {
				return nil
			},
			args: args{
				handler: func(ctx context.Context, command []Command) {},
			},
			wantErr: cqrs.ErrSecondArgumentOfHandlerMustBeStructOrPointerOfStruct,
		},
		{
			name: "handler has no return",
			prepare: func(commandBus cqrs.CommandBus) error {
				return nil
			},
			args: args{
				handler: func(ctx context.Context, command Command) {},
			},
			wantErr: cqrs.ErrHandlerMustHaveExactOneResult,
		},
		{
			name: "handler return bool",
			prepare: func(commandBus cqrs.CommandBus) error {
				return nil
			},
			args: args{
				handler: func(ctx context.Context, command Command) bool {
					return true
				},
			},
			wantErr: cqrs.ErrHandlerResultMustBeError,
		},
		{
			name: "handler already registered",
			prepare: func(commandBus cqrs.CommandBus) error {
				if err := commandBus.Register(func(ctx context.Context, command Command) error {
					return nil
				}); err != nil {
					return err
				}

				return nil
			},
			args: args{
				handler: func(ctx context.Context, command Command) error {
					return nil
				},
			},
			wantErr: cqrs.ErrCommandAlreadyRegistered,
		},
		{
			name: "success",
			prepare: func(commandBus cqrs.CommandBus) error {
				return nil
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
			commandBus := cqrs.NewCommandBus()

			err := tt.prepare(commandBus)
			if assert.NoError(t, err) {
				err := commandBus.Register(tt.args.handler)
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func Test_commandBus_Execute(t *testing.T) {
	t.Parallel()

	var (
		Err = errors.New("error")
	)

	type args struct {
		ctx     context.Context
		command Command
	}

	tests := []struct {
		name    string
		args    args
		prepare func(commandBus cqrs.CommandBus) error
		wantErr error
	}{
		{
			name: "command has not registered yet",
			args: args{
				ctx:     context.Background(),
				command: Command{},
			},
			prepare: func(commandBus cqrs.CommandBus) error {
				return nil
			},
			wantErr: cqrs.ErrCommandHasNotRegisteredYet,
		},
		{
			name: "command handler return error",
			args: args{
				ctx:     context.Background(),
				command: Command{},
			},
			prepare: func(commandBus cqrs.CommandBus) error {
				if err := commandBus.Register(func(ctx context.Context, command Command) error {
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
				ctx:     context.Background(),
				command: Command{},
			},
			prepare: func(commandBus cqrs.CommandBus) error {
				commandBus.Use(func(handler cqrs.CommandHandlerFunc[any]) cqrs.CommandHandlerFunc[any] {
					return func(ctx context.Context, command any) error {
						return Err
					}
				})

				if err := commandBus.Register(func(ctx context.Context, command Command) error {
					return nil
				}); err != nil {
					return err
				}

				return nil
			},
			wantErr: Err,
		},
		{
			name: "command handler return nil",
			args: args{
				ctx:     context.Background(),
				command: Command{},
			},
			prepare: func(commandBus cqrs.CommandBus) error {
				if err := commandBus.Register(func(ctx context.Context, command Command) error {
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
			commandBus := cqrs.NewCommandBus()

			if err := tt.prepare(commandBus); assert.NoError(t, err) {
				err = commandBus.Execute(tt.args.ctx, tt.args.command)
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
