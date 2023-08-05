package middlewares_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vulpes-ferrilata/cqrs"
	"github.com/vulpes-ferrilata/cqrs/middlewares"
	mock_mongo "github.com/vulpes-ferrilata/cqrs/pkg/mongo/mocks"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/mock/gomock"
)

func TestMongoTransactionMiddleware_CommandMiddleware(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.Background()
		command = struct{}{}

		Err = errors.New("error")
	)

	type mocks struct {
		db      *mock_mongo.MockDatabase
		client  *mock_mongo.MockClient
		session *mock_mongo.MockSession
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
			name: "start session fail",
			prepare: func(mocks mocks) {
				mocks.db.EXPECT().Client().Return(mocks.client)
				mocks.client.EXPECT().StartSession().Return(nil, Err)
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
				mocks.db.EXPECT().Client().Return(mocks.client)
				mocks.client.EXPECT().StartSession().Return(mocks.session, nil)
				mocks.session.EXPECT().WithTransaction(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error) {
					session := mongo.SessionFromContext(ctx)
					sessionCtx := mongo.NewSessionContext(ctx, session)
					return fn(sessionCtx)
				})
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
			name: "success",
			prepare: func(mocks mocks) {
				mocks.db.EXPECT().Client().Return(mocks.client)
				mocks.client.EXPECT().StartSession().Return(mocks.session, nil)
				mocks.session.EXPECT().WithTransaction(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error) {
					session := mongo.SessionFromContext(ctx)
					sessionCtx := mongo.NewSessionContext(ctx, session)
					return fn(sessionCtx)
				})
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
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mocks := mocks{
				db:      mock_mongo.NewMockDatabase(mockCtrl),
				client:  mock_mongo.NewMockClient(mockCtrl),
				session: mock_mongo.NewMockSession(mockCtrl),
			}

			tt.prepare(mocks)

			transactionMiddleware := middlewares.NewMongoTransactionMiddleware(mocks.db)
			commandMiddleware := transactionMiddleware.CommandMiddleware()
			handler := commandMiddleware(tt.args.handler)
			err := handler(tt.args.ctx, tt.args.command)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
