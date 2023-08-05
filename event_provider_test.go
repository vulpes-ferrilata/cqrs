package cqrs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vulpes-ferrilata/cqrs"
)

func Test_eventProvider_GetEvents(t *testing.T) {
	t.Parallel()

	var (
		event = Event{}
	)

	tests := []struct {
		name    string
		prepare func(eventProvider cqrs.EventProvider)
		want    []interface{}
	}{
		{
			name:    "initial",
			prepare: func(eventProvider cqrs.EventProvider) {},
			want:    nil,
		},
		{
			name: "event registered",
			prepare: func(eventProvider cqrs.EventProvider) {
				eventProvider.CollectEvents(event)
			},
			want: []interface{}{
				event,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventProvider := cqrs.NewEventProvider()

			tt.prepare(eventProvider)

			got := eventProvider.GetEvents()
			assert.Equal(t, tt.want, got)
		})
	}
}
