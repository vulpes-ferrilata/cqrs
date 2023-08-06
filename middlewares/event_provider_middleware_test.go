package middlewares_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vulpes-ferrilata/cqrs"
	"github.com/vulpes-ferrilata/cqrs/middlewares"
)

func TestEventProviderMiddleware_CommandMiddleware(t *testing.T) {
	t.Parallel()

	var (
		Err = errors.New("error")
	)

	type args struct {
		handler cqrs.CommandHandlerFunc[any]
		ctx     context.Context
		command interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "handler return error",
			args: args{
				handler: func(ctx context.Context, command any) error {
					_, ok := cqrs.GetEventProvider(ctx)
					if !ok {
						return cqrs.ErrEventProviderNotFound
					}

					return Err
				},
				ctx:     context.Background(),
				command: struct{}{},
			},
			wantErr: Err,
		},
		{
			name: "success",
			args: args{
				handler: func(ctx context.Context, command any) error {
					_, ok := cqrs.GetEventProvider(ctx)
					if !ok {
						return cqrs.ErrEventProviderNotFound
					}

					return nil
				},
				ctx:     context.Background(),
				command: struct{}{},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventProviderMiddleware := middlewares.NewEventProviderMiddleware()

			commandMiddleware := eventProviderMiddleware.CommandMiddleware()
			handler := commandMiddleware(tt.args.handler)
			err := handler(tt.args.ctx, tt.args.command)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
