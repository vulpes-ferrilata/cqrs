package cqrs_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vulpes-ferrilata/cqrs"
)

func TestGetEventProvider(t *testing.T) {
	t.Parallel()

	var (
		eventProvider = cqrs.NewEventProvider()
	)

	type wants struct {
		eventProvider cqrs.EventProvider
		ok            bool
	}
	tests := []struct {
		name    string
		prepare func() context.Context
		wants   wants
	}{
		{
			name: "no event provider injected into context",
			prepare: func() context.Context {
				return context.Background()
			},
			wants: wants{
				eventProvider: nil,
				ok:            false,
			},
		},
		{
			name: " event provider injected into context",
			prepare: func() context.Context {
				ctx := context.Background()
				ctx = cqrs.WithEventProvider(ctx, eventProvider)
				return ctx
			},
			wants: wants{
				eventProvider: eventProvider,
				ok:            true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.prepare()
			got, ok := cqrs.GetEventProvider(ctx)
			assert.Equal(t, tt.wants.eventProvider, got)
			assert.Equal(t, tt.wants.ok, ok)
		})
	}
}
