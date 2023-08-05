package mongo

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock_$GOPACKAGE
import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session interface {
	WithTransaction(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error)
}

func wrapSession(s mongo.Session) Session {
	return &session{
		Session: s,
	}
}

type session struct {
	mongo.Session
}
