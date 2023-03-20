package eventbus_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/vulpes-ferrilata/cqrs/eventbus"
)

type EventHandler struct{}

func (e EventHandler) Handle(ctx context.Context, event *Event) error {
	return nil
}

var _ = Describe("Generics", func() {
	var eventBus eventbus.EventBus

	BeforeEach(func() {
		eventBus = eventbus.NewEventBus()
	})

	Describe("RegisterEventHandlerWithEventBus", func() {
		var err error

		BeforeEach(func() {
			eventHandler := &EventHandler{}
			err = eventbus.RegisterEventHandlerWithEventBus[Event](eventBus, eventHandler)
		})

		It("should not return any error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("RegisterEventHandlerFuncWithEventBus", func() {
		var err error

		BeforeEach(func() {
			eventHandlerFunc := func(ctx context.Context, event *Event) error {
				return nil
			}
			err = eventbus.RegisterEventHandlerFuncWithEventBus[Event](eventBus, eventHandlerFunc)
		})

		It("should not return any error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
