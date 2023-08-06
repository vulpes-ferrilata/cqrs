package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type transactionKey struct{}

func withTransaction(ctx context.Context, transaction *mongo.Database) context.Context {
	return context.WithValue(ctx, transactionKey{}, transaction)
}

func getTransaction(ctx context.Context) (*mongo.Database, bool) {
	transaction, ok := ctx.Value(transactionKey{}).(*mongo.Database)
	return transaction, ok
}
