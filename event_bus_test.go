package cqrs_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/vulpes-ferrilata/cqrs"
)

type Event struct{}

var _ = Describe("EventBus", func() {
	var eventBus *cqrs.EventBus

	BeforeEach(func() {
		eventBus = &cqrs.EventBus{}
	})

	When("middleware was not registered", func() {
		When("event was not registered", func() {
			It("cannot publish event", func(ctx SpecContext) {
				err := eventBus.Publish(ctx, &Event{})
				Expect(err).Should(MatchError(cqrs.ErrHandlerNotFound))
			})
		})

		When("event was registered", func() {
			type Test struct {
				FirstHandlerExecuted  bool
				SecondHandlerExecuted bool
			}

			firstHandler := func(ctx context.Context, event *Event) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.FirstHandlerExecuted = true
				}

				return nil
			}

			secondHandler := func(ctx context.Context, event *Event) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.SecondHandlerExecuted = true
				}

				return nil
			}

			BeforeEach(func() {
				eventBus.Register(&Event{}, cqrs.WrapEventHandlerFunc(firstHandler))
				eventBus.Register(&Event{}, cqrs.WrapEventHandlerFunc(secondHandler))
			})

			It("can publish event", func(ctx SpecContext) {
				test := &Test{
					FirstHandlerExecuted:  false,
					SecondHandlerExecuted: false,
				}
				ctxWithTest := context.WithValue(ctx, "test", test)

				err := eventBus.Publish(ctxWithTest, &Event{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(test.FirstHandlerExecuted).Should(BeTrue())
				Expect(test.SecondHandlerExecuted).Should(BeTrue())
			})
		})
	})

	When("middleware was registered", func() {
		type Test struct {
			FirstHandlerExecuted  bool
			SecondHandlerExecuted bool
			Order                 []int
		}

		middlewareFunc1 := func(eventHandlerFunc cqrs.EventHandlerFunc[any]) cqrs.EventHandlerFunc[any] {
			return func(ctx context.Context, event any) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.Order = append(test.Order, 1)
				}

				return eventHandlerFunc(ctx, event)
			}
		}

		middlewareFunc2 := func(eventHandlerFunc cqrs.EventHandlerFunc[any]) cqrs.EventHandlerFunc[any] {
			return func(ctx context.Context, event any) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.Order = append(test.Order, 2)
				}

				return eventHandlerFunc(ctx, event)
			}
		}

		middlewareFunc3 := func(eventHandlerFunc cqrs.EventHandlerFunc[any]) cqrs.EventHandlerFunc[any] {
			return func(ctx context.Context, event any) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.Order = append(test.Order, 3)
				}

				return eventHandlerFunc(ctx, event)
			}
		}

		BeforeEach(func() {
			eventBus.Use(
				middlewareFunc1,
				middlewareFunc2,
				middlewareFunc3,
			)
		})

		When("event was not registered", func() {
			It("cannot publish event", func(ctx SpecContext) {
				err := eventBus.Publish(ctx, &Event{})
				Expect(err).Should(MatchError(cqrs.ErrHandlerNotFound))
			})
		})

		When("event was registered", func() {
			firstHandler := func(ctx context.Context, event *Event) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.FirstHandlerExecuted = true
					test.Order = append(test.Order, 4)
				}

				return nil
			}

			secondHandler := func(ctx context.Context, event *Event) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.SecondHandlerExecuted = true
					test.Order = append(test.Order, 4)
				}

				return nil
			}

			BeforeEach(func() {
				eventBus.Register(&Event{}, cqrs.WrapEventHandlerFunc(firstHandler))
				eventBus.Register(&Event{}, cqrs.WrapEventHandlerFunc(secondHandler))
			})

			It("can publish event", func(ctx SpecContext) {
				test := &Test{
					FirstHandlerExecuted:  false,
					SecondHandlerExecuted: false,
					Order:                 make([]int, 0),
				}
				ctxWithTest := context.WithValue(ctx, "test", test)

				err := eventBus.Publish(ctxWithTest, &Event{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(test.FirstHandlerExecuted).Should(BeTrue())
				Expect(test.SecondHandlerExecuted).Should(BeTrue())
				Expect(test.Order).Should(BeEquivalentTo([]int{1, 2, 3, 4, 1, 2, 3, 4}))
			})
		})
	})
})
