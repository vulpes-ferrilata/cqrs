package gorm

import (
	"context"

	"gorm.io/gorm"
)

type transactionKey struct{}

func withTransaction(ctx context.Context, transaction *gorm.DB) context.Context {
	return context.WithValue(ctx, transactionKey{}, transaction)
}

func getTransaction(ctx context.Context) (*gorm.DB, bool) {
	transaction, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	return transaction, ok
}
