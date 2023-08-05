package cqrs

import (
	"context"
)

type eventProviderKey struct{}

func WithEventProvider(ctx context.Context, eventProvider EventProvider) context.Context {
	return context.WithValue(ctx, eventProviderKey{}, eventProvider)
}

func GetEventProvider(ctx context.Context) (EventProvider, bool) {
	eventProvider, ok := ctx.Value(eventProviderKey{}).(EventProvider)
	return eventProvider, ok
}
