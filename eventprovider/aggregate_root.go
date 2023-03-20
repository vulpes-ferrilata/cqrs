package eventprovider

type AggregateRoot interface {
	GetEvents() []interface{}
}
