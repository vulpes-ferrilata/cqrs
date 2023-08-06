package db

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock_$GOPACKAGE
import "context"

type TransactionManager[DB any] interface {
	IsTransactionStarted(ctx context.Context) bool
	StartTransaction(ctx context.Context) (Committer, context.Context, error)
	GetTransaction(ctx context.Context) DB
}
