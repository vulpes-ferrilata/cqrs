package mongo

//go:generate mockgen -destination=./mocks/mock_$GOFILE -source=$GOFILE -package=mock_$GOPACKAGE
import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	StartSession(opts ...*options.SessionOptions) (Session, error)
}

func wrapClient(c *mongo.Client) Client {
	return &client{
		Client: c,
	}
}

type client struct {
	*mongo.Client
}

func (c client) StartSession(opts ...*options.SessionOptions) (Session, error) {
	session, err := c.Client.StartSession(opts...)
	return wrapSession(session), err
}
