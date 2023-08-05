package cqrs

import "sync"

type EventProvider interface {
	GetEvents() []interface{}
	CollectEvents(events ...interface{})
}

func NewEventProvider() EventProvider {
	return &eventProvider{}
}

type eventProvider struct {
	mu     sync.RWMutex
	events []interface{}
}

func (e *eventProvider) GetEvents() []interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.events
}

func (e *eventProvider) CollectEvents(events ...interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.events = append(e.events, events...)
}
