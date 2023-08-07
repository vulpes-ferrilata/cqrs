package gorm

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"github.com/vulpes-ferrilata/cqrs/pkg/db"
)

func NewTransactionManager(db *gorm.DB, opts *sql.TxOptions) db.TransactionManager[*gorm.DB] {
	return &transactionManager{
		db:   db,
		opts: opts,
	}
}

type transactionManager struct {
	db   *gorm.DB
	opts *sql.TxOptions
}

func (t transactionManager) IsTransactionStarted(ctx context.Context) bool {
	_, ok := getTransaction(ctx)
	return ok
}

func (t transactionManager) StartTransaction(ctx context.Context) (db.Committer, context.Context, error) {
	transaction := t.db.WithContext(ctx).Begin(t.opts)

	ctx = withTransaction(ctx, transaction)

	return newCommitter(transaction), ctx, nil
}

func (t transactionManager) GetTransaction(ctx context.Context) *gorm.DB {
	transaction, ok := getTransaction(ctx)
	if !ok {
		return t.db.WithContext(ctx)
	}

	return transaction
}
