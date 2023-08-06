package middlewares_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/vulpes-ferrilata/cqrs"
	"github.com/vulpes-ferrilata/cqrs/middlewares"
	mock_db "github.com/vulpes-ferrilata/cqrs/pkg/db/mocks"
)

func TestTransactionMiddleware_CommandMiddleware(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.Background()
		newCtx  = context.WithValue(ctx, "xxx", "yyy")
		command = struct{}{}

		Err = errors.New("error")
	)

	type mocks struct {
		transactionManager *mock_db.MockTransactionManager[*gorm.DB]
		committer          *mock_db.MockCommitter
	}
	type args struct {
		handler cqrs.CommandHandlerFunc[any]
		ctx     context.Context
		command interface{}
	}
	tests := []struct {
		name    string
		prepare func(mocks mocks)
		args    args
		wantErr error
	}{
		{
			name: "start transaction fail",
			prepare: func(mocks mocks) {
				mocks.transactionManager.EXPECT().IsTransactionStarted(ctx).Return(false)
				mocks.transactionManager.EXPECT().StartTransaction(ctx).Return(nil, nil, Err)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return nil
				},
				ctx:     ctx,
				command: command,
			},
			wantErr: Err,
		},
		{
			name: "handler return error",
			prepare: func(mocks mocks) {
				mocks.transactionManager.EXPECT().IsTransactionStarted(ctx).Return(false)
				mocks.transactionManager.EXPECT().StartTransaction(ctx).Return(mocks.committer, newCtx, nil)
				mocks.committer.EXPECT().RollbackTransaction(newCtx).Return(nil)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return Err
				},
				ctx:     ctx,
				command: command,
			},
			wantErr: Err,
		},
		{
			name: "handler panic",
			prepare: func(mocks mocks) {
				mocks.transactionManager.EXPECT().IsTransactionStarted(ctx).Return(false)
				mocks.transactionManager.EXPECT().StartTransaction(ctx).Return(mocks.committer, newCtx, nil)
				mocks.committer.EXPECT().RollbackTransaction(newCtx).Return(nil)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					panic(Err)
				},
				ctx:     ctx,
				command: command,
			},
			wantErr: Err,
		},
		{
			name: "commit transaction fail",
			prepare: func(mocks mocks) {
				mocks.transactionManager.EXPECT().IsTransactionStarted(ctx).Return(false)
				mocks.transactionManager.EXPECT().StartTransaction(ctx).Return(mocks.committer, newCtx, nil)
				mocks.committer.EXPECT().CommitTransaction(newCtx).Return(Err)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return nil
				},
				ctx:     ctx,
				command: command,
			},
			wantErr: Err,
		},
		{
			name: "commit transaction fail",
			prepare: func(mocks mocks) {
				mocks.transactionManager.EXPECT().IsTransactionStarted(ctx).Return(false)
				mocks.transactionManager.EXPECT().StartTransaction(ctx).Return(mocks.committer, newCtx, nil)
				mocks.committer.EXPECT().CommitTransaction(newCtx).Return(Err)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return nil
				},
				ctx:     ctx,
				command: command,
			},
			wantErr: Err,
		},
		{
			name: "success",
			prepare: func(mocks mocks) {
				mocks.transactionManager.EXPECT().IsTransactionStarted(ctx).Return(false)
				mocks.transactionManager.EXPECT().StartTransaction(ctx).Return(mocks.committer, newCtx, nil)
				mocks.committer.EXPECT().CommitTransaction(newCtx).Return(nil)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return nil
				},
				ctx:     ctx,
				command: command,
			},
			wantErr: nil,
		},
		{
			name: "transaction already started - handler return error",
			prepare: func(mocks mocks) {
				mocks.transactionManager.EXPECT().IsTransactionStarted(ctx).Return(true)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return Err
				},
				ctx:     ctx,
				command: command,
			},
			wantErr: Err,
		},
		{
			name: "transaction already started - success",
			prepare: func(mocks mocks) {
				mocks.transactionManager.EXPECT().IsTransactionStarted(ctx).Return(true)
			},
			args: args{
				handler: func(ctx context.Context, command interface{}) error {
					return nil
				},
				ctx:     ctx,
				command: command,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				recover()
			}()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mocks := mocks{
				transactionManager: mock_db.NewMockTransactionManager[*gorm.DB](mockCtrl),
				committer:          mock_db.NewMockCommitter(mockCtrl),
			}

			tt.prepare(mocks)

			transactionMiddleware := middlewares.NewTransactionMiddleware[*gorm.DB](mocks.transactionManager)
			commandMiddleware := transactionMiddleware.CommandMiddleware()
			handler := commandMiddleware(tt.args.handler)
			err := handler(tt.args.ctx, tt.args.command)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
