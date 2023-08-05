package mongo

import "go.mongodb.org/mongo-driver/mongo"

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock_$GOPACKAGE
type Database interface {
	Client() Client
}

func WrapDatabase(db *mongo.Database) Database {
	return &database{
		Database: db,
	}
}

type database struct {
	*mongo.Database
}

func (d database) Client() Client {
	return wrapClient(d.Database.Client())
}
