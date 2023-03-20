package eventprovider

import "sync"

type EventProvider struct {
	mu     sync.RWMutex
	events []interface{}
}

func (e *EventProvider) GetEvents() []interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.events
}

func (e *EventProvider) CollectEvents(aggregateRoot AggregateRoot) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.events = append(e.events, aggregateRoot.GetEvents()...)
}
