package db

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock_$GOPACKAGE
import "context"

type Committer interface {
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context) error
}
