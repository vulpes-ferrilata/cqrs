package eventbus_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/vulpes-ferrilata/cqrs/v2/eventbus"
)

type EventHandler struct{}

func (e EventHandler) Handle(ctx context.Context, event *Event) error {
	return nil
}

type EventBusMock struct {
	mock.Mock
}

func (e EventBusMock) Use(middlewareFunc eventbus.EventMiddlewareFunc) error {
	args := e.Called(middlewareFunc)
	return args.Error(0)
}

func (e EventBusMock) Register(handlerFunc interface{}) error {
	args := e.Called(handlerFunc)
	return args.Error(0)
}

func (e EventBusMock) Publish(ctx context.Context, events ...interface{}) error {
	args := e.Called(ctx, events)
	return args.Error(0)
}

var _ = Describe("Generics", func() {
	var eventBus *EventBusMock

	BeforeEach(func() {
		eventBus = &EventBusMock{}
	})

	Describe("RegisterEventHandlerWithEventBus", func() {
		Context("when EventBus return error", func() {
			var err error

			expectedErr := errors.New("error")

			BeforeEach(func() {
				eventHandler := &EventHandler{}

				eventBus.
					On("Register", mock.AnythingOfType("func(context.Context, *eventbus_test.Event) error")).
					Return(expectedErr).
					Once()

				err = eventbus.RegisterEventHandlerWithEventBus[Event](eventBus, eventHandler)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(expectedErr))
			})
		})

		Context("when EventBus not return any error", func() {
			var err error

			BeforeEach(func() {
				eventHandler := &EventHandler{}

				eventBus.
					On("Register", mock.AnythingOfType("func(context.Context, *eventbus_test.Event) error")).
					Return(nil).
					Once()

				err = eventbus.RegisterEventHandlerWithEventBus[Event](eventBus, eventHandler)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("RegisterEventHandlerFuncWithEventBus", func() {
		Context("when EventBus return error", func() {
			var err error

			expectedErr := errors.New("error")

			BeforeEach(func() {
				eventHandlerFunc := func(ctx context.Context, event *Event) error {
					return nil
				}

				eventBus.
					On("Register", mock.AnythingOfType("EventHandlerFunc[*github.com/vulpes-ferrilata/cqrs/v2/eventbus_test.Event]")).
					Return(expectedErr).
					Once()

				err = eventbus.RegisterEventHandlerFuncWithEventBus[Event](eventBus, eventHandlerFunc)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(expectedErr))
			})
		})

		Context("when EventBus not return any error", func() {
			var err error

			BeforeEach(func() {
				eventHandlerFunc := func(ctx context.Context, event *Event) error {
					return nil
				}

				eventBus.
					On("Register", mock.AnythingOfType("EventHandlerFunc[*github.com/vulpes-ferrilata/cqrs/v2/eventbus_test.Event]")).
					Return(nil).
					Once()

				err = eventbus.RegisterEventHandlerFuncWithEventBus[Event](eventBus, eventHandlerFunc)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
