package middlewares

import (
	"context"

	"github.com/vulpes-ferrilata/cqrs"
	"github.com/vulpes-ferrilata/cqrs/pkg/db"
)

func NewTransactionMiddleware[DB any](transactionManager db.TransactionManager[DB]) *TransactionMiddleware[DB] {
	return &TransactionMiddleware[DB]{
		transactionManager: transactionManager,
	}
}

type TransactionMiddleware[DB any] struct {
	transactionManager db.TransactionManager[DB]
}

func (m TransactionMiddleware[DB]) CommandMiddleware() cqrs.CommandMiddlewareFunc {
	return func(handler cqrs.CommandHandlerFunc[any]) cqrs.CommandHandlerFunc[any] {
		return func(ctx context.Context, command any) error {
			if isTransactionStarted := m.transactionManager.IsTransactionStarted(ctx); !isTransactionStarted {
				committer, ctx, err := m.transactionManager.StartTransaction(ctx)
				if err != nil {
					return err
				}
				defer func() {
					if r := recover(); r != nil {
						committer.RollbackTransaction(ctx)

						panic(r)
					}
				}()

				if err := handler(ctx, command); err != nil {
					committer.RollbackTransaction(ctx)

					return err
				}

				if err := committer.CommitTransaction(ctx); err != nil {
					return err
				}

				return nil
			}

			if err := handler(ctx, command); err != nil {
				return err
			}

			return nil
		}
	}
}
