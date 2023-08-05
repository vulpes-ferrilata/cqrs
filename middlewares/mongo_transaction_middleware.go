package middlewares

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vulpes-ferrilata/cqrs"
	pkg_mongo "github.com/vulpes-ferrilata/cqrs/pkg/mongo"
)

func NewMongoTransactionMiddleware(db pkg_mongo.Database, opts ...*options.SessionOptions) *MongoTransactionMiddleware {
	return &MongoTransactionMiddleware{
		db:   db,
		opts: opts,
	}
}

type MongoTransactionMiddleware struct {
	db   pkg_mongo.Database
	opts []*options.SessionOptions
}

func (m MongoTransactionMiddleware) CommandMiddleware() cqrs.CommandMiddlewareFunc {
	return func(handler cqrs.CommandHandlerFunc[any]) cqrs.CommandHandlerFunc[any] {
		return func(ctx context.Context, command any) error {
			session, err := m.db.Client().StartSession(m.opts...)
			if err != nil {
				return err
			}

			if _, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
				if err := handler(ctx, command); err != nil {
					return nil, err
				}

				return nil, nil
			}); err != nil {
				return err
			}

			return nil
		}
	}
}
