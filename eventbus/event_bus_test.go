package eventbus_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vulpes-ferrilata/cqrs/eventbus"
)

type HandlerType string

const (
	FirstMiddlewareBeforeHandler  HandlerType = "first middleware before handler"
	SecondMiddlewareBeforeHandler HandlerType = "second middleware before handler"
	ThirdMiddlewareBeforeHandler  HandlerType = "third middleware before handler"
	Handler                       HandlerType = "handler"
	ThirdMiddlewareAfterHandler   HandlerType = "third middleware after handler"
	SecondMiddlewareAfterHandler  HandlerType = "second middleware after handler"
	FirstMiddlewareAfterHandler   HandlerType = "first middleware after handler"
)

type Event struct{}

var _ = Describe("EventBus", func() {
	var ctx context.Context
	var eventBus eventbus.EventBus

	BeforeEach(func() {
		ctx = context.Background()
		eventBus = eventbus.NewEventBus()
	})

	Describe("Use", func() {
		Context("when middlewareFunc is nil", func() {
			var err error

			BeforeEach(func() {
				err = eventBus.Use(nil)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(eventbus.ErrMiddlewareFuncMustNotBeNil))
			})
		})

		Context("when middlewareFunc is function", func() {
			var err error

			BeforeEach(func() {
				middlewareFunc := func(eventHandlerFunc eventbus.EventHandlerFunc[any]) eventbus.EventHandlerFunc[any] {
					return func(ctx context.Context, event interface{}) error {
						return eventHandlerFunc(ctx, event)
					}
				}

				err = eventBus.Use(middlewareFunc)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("Register", func() {
		DescribeTable("when handlerFunc is invalid",
			func(handlerFunc interface{}, expectedErr error) {
				err := eventBus.Register(handlerFunc)
				Expect(err).Should(MatchError(expectedErr))
			},
			Entry("when handlerFunc is nil", nil, eventbus.ErrHandlerFuncMustBeNonNilFunction),
			Entry("when handlerFunc is not function", struct{}{}, eventbus.ErrHandlerFuncMustBeNonNilFunction),
			Entry("when handlerFunc has 1 argument", func(context.Context) {}, eventbus.ErrHandlerFuncMustHaveExactTwoArguments),
			Entry("when handlerFunc has 3 arguments", func(context.Context, *Event, bool) {}, eventbus.ErrHandlerFuncMustHaveExactTwoArguments),
			Entry("when first argument of handlerFunc is not context.Context", func(struct{}, *Event) {}, eventbus.ErrFirstArgumentOfHandlerMustBeContext),
			Entry("when second argument of handlerFunc is integer", func(context.Context, int) {}, eventbus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when second argument of handlerFunc is struct", func(context.Context, Event) {}, eventbus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when second argument of handlerFunc is slice of pointer of struct", func(context.Context, []Event) {}, eventbus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when second argument of handlerFunc is pointer of pointer of struct", func(context.Context, **Event) {}, eventbus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when handlerFunc has no result", func(context.Context, *Event) {}, eventbus.ErrHandlerFuncMustHaveExactOneResult),
			Entry("when handlerFunc has 2 results", func(context.Context, *Event) (interface{}, error) { return nil, nil }, eventbus.ErrHandlerFuncMustHaveExactOneResult),
			Entry("when result of handlerFunc is not error", func(context.Context, *Event) interface{} { return nil }, eventbus.ErrResultMustBeError),
		)

		Context("when handlerFunc is valid", func() {
			var err error

			BeforeEach(func() {
				handlerFunc := func(context.Context, *Event) error {
					return nil
				}

				err = eventBus.Register(handlerFunc)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when register multiple handlerFunc with same event", func() {
			var err error

			BeforeEach(func() {
				for i := 1; i <= 2; i++ {
					handlerFunc := func(context.Context, *Event) error {
						return nil
					}

					err = eventBus.Register(handlerFunc)
				}
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("Execute", func() {
		Context("when handler has not registered yet", func() {
			Context("when middleware has not registered yet", func() {
				DescribeTable("when event is invalid",
					func(event interface{}, expectedErr error) {
						err := eventBus.Publish(ctx, event)
						Expect(err).Should(MatchError(expectedErr))
					},
					Entry("when event is nil", nil, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is integer", 5, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is struct", Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is slice of struct", []Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is pointer of pointer of struct", Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
				)

				Context("when event is valid", func() {
					var err error

					BeforeEach(func() {
						event := &Event{}
						err = eventBus.Publish(ctx, event)
					})

					It("should not return any error", func() {
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
		})

		Context("when handler is registered", func() {
			var handlerExecuted bool
			var executedOrders []HandlerType

			BeforeEach(func() {
				handlerExecuted = false
				executedOrders = make([]HandlerType, 0)

				handlerFunc := func(context.Context, *Event) error {
					handlerExecuted = true
					executedOrders = append(executedOrders, Handler)
					return nil
				}

				err := eventBus.Register(handlerFunc)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Context("when middleware has not registered yet", func() {
				DescribeTable("when event is invalid",
					func(event interface{}, expectedErr error) {
						err := eventBus.Publish(ctx, event)
						Expect(err).Should(MatchError(expectedErr))
					},
					Entry("when event is nil", nil, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is integer", 5, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is struct", Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is slice of struct", []Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is pointer of pointer of struct", Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
				)

				Context("when event is valid", func() {
					var err error

					BeforeEach(func() {
						event := &Event{}
						err = eventBus.Publish(ctx, event)
					})

					It("should not return any error", func() {
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should execute handler", func() {
						Expect(handlerExecuted).Should(BeTrue())
					})

					It("should have correct execution orders", func() {
						Expect(executedOrders).Should(HaveExactElements(Handler))
					})
				})
			})

			Context("when middleware is registered", func() {
				BeforeEach(func() {
					firstMiddleware := func(eventHandlerFunc eventbus.EventHandlerFunc[any]) eventbus.EventHandlerFunc[any] {
						return func(ctx context.Context, event any) error {
							executedOrders = append(executedOrders, FirstMiddlewareBeforeHandler)
							err := eventHandlerFunc(ctx, event)
							executedOrders = append(executedOrders, FirstMiddlewareAfterHandler)
							return err
						}
					}

					secondMiddleware := func(eventHandlerFunc eventbus.EventHandlerFunc[any]) eventbus.EventHandlerFunc[any] {
						return func(ctx context.Context, event any) error {
							executedOrders = append(executedOrders, SecondMiddlewareBeforeHandler)
							err := eventHandlerFunc(ctx, event)
							executedOrders = append(executedOrders, SecondMiddlewareAfterHandler)
							return err
						}
					}

					thirdMiddleware := func(eventHandlerFunc eventbus.EventHandlerFunc[any]) eventbus.EventHandlerFunc[any] {
						return func(ctx context.Context, event any) error {
							executedOrders = append(executedOrders, ThirdMiddlewareBeforeHandler)
							err := eventHandlerFunc(ctx, event)
							executedOrders = append(executedOrders, ThirdMiddlewareAfterHandler)
							return err
						}
					}

					err := eventBus.Use(firstMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = eventBus.Use(secondMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = eventBus.Use(thirdMiddleware)
					Expect(err).ShouldNot(HaveOccurred())
				})

				DescribeTable("when event is invalid",
					func(event interface{}, expectedErr error) {
						err := eventBus.Publish(ctx, event)
						Expect(err).Should(MatchError(expectedErr))
					},
					Entry("when event is nil", nil, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is integer", 5, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is struct", Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is slice of struct", []Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is pointer of pointer of struct", Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
				)

				Context("when event is valid", func() {
					var err error

					BeforeEach(func() {
						event := &Event{}
						err = eventBus.Publish(ctx, event)
					})

					It("should not return any error", func() {
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should execute handler", func() {
						Expect(handlerExecuted).Should(BeTrue())
					})

					It("should have correct execution orders", func() {
						Expect(executedOrders).Should(
							HaveExactElements(
								FirstMiddlewareBeforeHandler,
								SecondMiddlewareBeforeHandler,
								ThirdMiddlewareBeforeHandler,
								Handler,
								ThirdMiddlewareAfterHandler,
								SecondMiddlewareAfterHandler,
								FirstMiddlewareAfterHandler,
							),
						)
					})
				})
			})
		})

		Context("when handler with error return is registered", func() {
			var handlerExecuted bool
			var executedOrders []HandlerType
			Err := errors.New("error")

			BeforeEach(func() {
				handlerExecuted = false
				executedOrders = make([]HandlerType, 0)

				handlerFunc := func(context.Context, *Event) error {
					handlerExecuted = true
					executedOrders = append(executedOrders, Handler)
					return Err
				}

				err := eventBus.Register(handlerFunc)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Context("when middleware has not registered yet", func() {
				DescribeTable("when event is invalid",
					func(event interface{}, expectedErr error) {
						err := eventBus.Publish(ctx, event)
						Expect(err).Should(MatchError(expectedErr))
					},
					Entry("when event is nil", nil, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is integer", 5, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is struct", Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is slice of struct", []Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is pointer of pointer of struct", Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
				)

				Context("when event is valid", func() {
					var err error

					BeforeEach(func() {
						event := &Event{}
						err = eventBus.Publish(ctx, event)
					})

					It("should return error", func() {
						Expect(err).Should(MatchError(Err))
					})

					It("should execute handler", func() {
						Expect(handlerExecuted).Should(BeTrue())
					})

					It("should have correct execution orders", func() {
						Expect(executedOrders).Should(HaveExactElements(Handler))
					})
				})
			})

			Context("when middleware is registered", func() {
				BeforeEach(func() {
					firstMiddleware := func(eventHandlerFunc eventbus.EventHandlerFunc[any]) eventbus.EventHandlerFunc[any] {
						return func(ctx context.Context, event any) error {
							executedOrders = append(executedOrders, FirstMiddlewareBeforeHandler)
							err := eventHandlerFunc(ctx, event)
							executedOrders = append(executedOrders, FirstMiddlewareAfterHandler)
							return err
						}
					}

					secondMiddleware := func(eventHandlerFunc eventbus.EventHandlerFunc[any]) eventbus.EventHandlerFunc[any] {
						return func(ctx context.Context, event any) error {
							executedOrders = append(executedOrders, SecondMiddlewareBeforeHandler)
							err := eventHandlerFunc(ctx, event)
							executedOrders = append(executedOrders, SecondMiddlewareAfterHandler)
							return err
						}
					}

					thirdMiddleware := func(eventHandlerFunc eventbus.EventHandlerFunc[any]) eventbus.EventHandlerFunc[any] {
						return func(ctx context.Context, event any) error {
							executedOrders = append(executedOrders, ThirdMiddlewareBeforeHandler)
							err := eventHandlerFunc(ctx, event)
							executedOrders = append(executedOrders, ThirdMiddlewareAfterHandler)
							return err
						}
					}

					err := eventBus.Use(firstMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = eventBus.Use(secondMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = eventBus.Use(thirdMiddleware)
					Expect(err).ShouldNot(HaveOccurred())
				})

				DescribeTable("when event is invalid",
					func(event interface{}, expectedErr error) {
						err := eventBus.Publish(ctx, event)
						Expect(err).Should(MatchError(expectedErr))
					},
					Entry("when event is nil", nil, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is integer", 5, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is struct", Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is slice of struct", []Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
					Entry("when event is pointer of pointer of struct", Event{}, eventbus.ErrEventMustBeNonNilPointerOfStruct),
				)

				Context("when event is valid", func() {
					var err error

					BeforeEach(func() {
						event := &Event{}
						err = eventBus.Publish(ctx, event)
					})

					It("should return error", func() {
						Expect(err).Should(MatchError(Err))
					})

					It("should execute handler", func() {
						Expect(handlerExecuted).Should(BeTrue())
					})

					It("should have correct execution orders", func() {
						Expect(executedOrders).Should(
							HaveExactElements(
								FirstMiddlewareBeforeHandler,
								SecondMiddlewareBeforeHandler,
								ThirdMiddlewareBeforeHandler,
								Handler,
								ThirdMiddlewareAfterHandler,
								SecondMiddlewareAfterHandler,
								FirstMiddlewareAfterHandler,
							),
						)
					})
				})
			})
		})
	})
})
