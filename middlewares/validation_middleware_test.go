package middlewares_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/vulpes-ferrilata/cqrs"
	"github.com/vulpes-ferrilata/cqrs/middlewares"
	mock_validator "github.com/vulpes-ferrilata/cqrs/pkg/validator/mocks"
)

func TestValidationMiddleware_CommandMiddleware(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.Background()
		command = struct{}{}

		Err = errors.New("error")
	)

	type fields struct {
		validate *mock_validator.MockValidate
	}
	type args struct {
		handler cqrs.CommandHandlerFunc[any]
		ctx     context.Context
		command interface{}
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name    string
		prepare func(fields fields)
		args    args
		wants   wants
	}{
		{
			name: "validator return error",
			prepare: func(fields fields) {
				fields.validate.EXPECT().StructCtx(ctx, command).Return(Err)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return nil
				},
				ctx:     ctx,
				command: command,
			},
			wants: wants{
				err: Err,
			},
		},
		{
			name: "handler return error",
			prepare: func(fields fields) {
				fields.validate.EXPECT().StructCtx(ctx, command).Return(nil)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return Err
				},
				ctx:     ctx,
				command: command,
			},
			wants: wants{
				err: Err,
			},
		},
		{
			name: "success",
			prepare: func(fields fields) {
				fields.validate.EXPECT().StructCtx(ctx, command).Return(nil)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return nil
				},
				ctx:     ctx,
				command: command,
			},
			wants: wants{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fields := fields{
				validate: mock_validator.NewMockValidate(mockCtrl),
			}

			tt.prepare(fields)

			validationMiddleware := middlewares.NewValidationMiddleware(fields.validate)
			commandMiddleware := validationMiddleware.CommandMiddleware()
			handler := commandMiddleware(tt.args.handler)
			err := handler(ctx, command)
			assert.ErrorIs(t, err, tt.wants.err)
		})
	}
}

func TestValidationMiddleware_QueryMiddleware(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.Background()
		command = struct{}{}
		result  = struct{}{}
		Err     = errors.New("error")
	)

	type mocks struct {
		validate *mock_validator.MockValidate
	}
	type args struct {
		handler cqrs.QueryHandlerFunc[any, any]
		ctx     context.Context
		command interface{}
	}
	type wants struct {
		result interface{}
		err    error
	}
	tests := []struct {
		name    string
		prepare func(mocks mocks)
		args    args
		wants   wants
	}{
		{
			name: "validator return error",
			prepare: func(mocks mocks) {
				mocks.validate.EXPECT().StructCtx(ctx, command).Return(Err)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) (interface{}, error) {
					return result, nil
				},
				ctx:     ctx,
				command: command,
			},
			wants: wants{
				result: nil,
				err:    Err,
			},
		},
		{
			name: "handler return error",
			prepare: func(mocks mocks) {
				mocks.validate.EXPECT().StructCtx(ctx, command).Return(nil)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) (interface{}, error) {
					return nil, Err
				},
				ctx:     ctx,
				command: command,
			},
			wants: wants{
				result: nil,
				err:    Err,
			},
		},
		{
			name: "success",
			prepare: func(mocks mocks) {
				mocks.validate.EXPECT().StructCtx(ctx, command).Return(nil)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) (interface{}, error) {
					return result, nil
				},
				ctx:     ctx,
				command: command,
			},
			wants: wants{
				result: result,
				err:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mocks := mocks{
				validate: mock_validator.NewMockValidate(mockCtrl),
			}

			tt.prepare(mocks)

			validationMiddleware := middlewares.NewValidationMiddleware(mocks.validate)
			queryMiddleware := validationMiddleware.QueryMiddleware()
			handler := queryMiddleware(tt.args.handler)
			result, err := handler(ctx, command)
			assert.ErrorIs(t, err, tt.wants.err)
			assert.Equal(t, tt.wants.result, result)
		})
	}
}
