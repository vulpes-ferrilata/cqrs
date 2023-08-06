package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vulpes-ferrilata/cqrs/pkg/db"
)

func newCommitter(session mongo.Session) db.Committer {
	return &committer{
		session: session,
	}
}

type committer struct {
	session mongo.Session
}

func (c committer) CommitTransaction(ctx context.Context) error {
	return c.session.CommitTransaction(ctx)
}

func (c committer) RollbackTransaction(ctx context.Context) error {
	return c.session.AbortTransaction(ctx)
}
