package context

import (
	"context"

	"github.com/vulpes-ferrilata/cqrs/eventprovider"
)

type eventProviderKey struct{}

func WithEventProvider(ctx context.Context, eventProvider *eventprovider.EventProvider) context.Context {
	return context.WithValue(ctx, eventProviderKey{}, eventProvider)
}

func GetEventProvider(ctx context.Context) (*eventprovider.EventProvider, error) {
	eventProvider, ok := ctx.Value(eventProviderKey{}).(*eventprovider.EventProvider)
	if !ok {
		return nil, ErrEventProviderNotFound
	}

	return eventProvider, nil
}
