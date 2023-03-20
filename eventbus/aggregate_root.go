package eventbus

type AggregateRoot interface {
	GetEvents() []interface{}
}
