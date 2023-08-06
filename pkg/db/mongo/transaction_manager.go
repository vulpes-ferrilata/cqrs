package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vulpes-ferrilata/cqrs/pkg/db"
)

func NewTransactionManager(db *mongo.Database,
	sessionOptions *options.SessionOptions,
	transactionOptions *options.TransactionOptions) db.TransactionManager[*mongo.Database] {
	return &transactionManager{
		db:                 db,
		sessionOptions:     sessionOptions,
		transactionOptions: transactionOptions,
	}
}

type transactionManager struct {
	db                 *mongo.Database
	sessionOptions     *options.SessionOptions
	transactionOptions *options.TransactionOptions
}

func (t transactionManager) IsTransactionStarted(ctx context.Context) bool {
	_, ok := getTransaction(ctx)
	return ok
}

func (t transactionManager) StartTransaction(ctx context.Context) (db.Committer, context.Context, error) {
	session, err := t.db.Client().StartSession(t.sessionOptions)
	if err != nil {
		return nil, ctx, err
	}

	if err = session.StartTransaction(t.transactionOptions); err != nil {
		return nil, ctx, err
	}

	ctx = withTransaction(ctx, t.db)

	return newCommitter(session), ctx, nil
}

func (t transactionManager) GetTransaction(ctx context.Context) *mongo.Database {
	transaction, ok := getTransaction(ctx)
	if !ok {
		return t.db
	}

	return transaction
}
