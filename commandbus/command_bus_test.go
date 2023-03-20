package commandbus_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vulpes-ferrilata/cqrs/commandbus"
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

type Command struct{}

var _ = Describe("CommandBus", func() {
	var ctx context.Context
	var commandBus commandbus.CommandBus

	BeforeEach(func() {
		ctx = context.Background()
		commandBus = commandbus.NewCommandBus()
	})

	Describe("Use", func() {
		Context("when middlewareFunc is nil", func() {
			var err error

			BeforeEach(func() {
				err = commandBus.Use(nil)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(commandbus.ErrMiddlewareFuncMustNotBeNil))
			})
		})

		Context("when middlewareFunc is function", func() {
			var err error

			BeforeEach(func() {
				middlewareFunc := func(commandHandlerFunc commandbus.CommandHandlerFunc[any]) commandbus.CommandHandlerFunc[any] {
					return func(ctx context.Context, command interface{}) error {
						return commandHandlerFunc(ctx, command)
					}
				}

				err = commandBus.Use(middlewareFunc)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("Register", func() {
		DescribeTable("when handlerFunc is invalid",
			func(handlerFunc interface{}, expectedErr error) {
				err := commandBus.Register(handlerFunc)
				Expect(err).Should(MatchError(expectedErr))
			},
			Entry("when handlerFunc is nil", nil, commandbus.ErrHandlerFuncMustBeNonNilFunction),
			Entry("when handlerFunc is not function", struct{}{}, commandbus.ErrHandlerFuncMustBeNonNilFunction),
			Entry("when handlerFunc has 1 argument", func(context.Context) {}, commandbus.ErrHandlerFuncMustHaveExactTwoArguments),
			Entry("when handlerFunc has 3 arguments", func(context.Context, *Command, bool) {}, commandbus.ErrHandlerFuncMustHaveExactTwoArguments),
			Entry("when first argument of handlerFunc is not context.Context", func(struct{}, *Command) {}, commandbus.ErrFirstArgumentOfHandlerMustBeContext),
			Entry("when second argument of handlerFunc is integer", func(context.Context, int) {}, commandbus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when second argument of handlerFunc is struct", func(context.Context, Command) {}, commandbus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when second argument of handlerFunc is slice of pointer of struct", func(context.Context, []Command) {}, commandbus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when second argument of handlerFunc is pointer of pointer of struct", func(context.Context, **Command) {}, commandbus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when handlerFunc has no result", func(context.Context, *Command) {}, commandbus.ErrHandlerFuncMustHaveExactOneResult),
			Entry("when handlerFunc has 2 results", func(context.Context, *Command) (interface{}, error) { return nil, nil }, commandbus.ErrHandlerFuncMustHaveExactOneResult),
			Entry("when result of handlerFunc is not error", func(context.Context, *Command) interface{} { return nil }, commandbus.ErrResultMustBeError),
		)

		Context("when handlerFunc is valid", func() {
			var err error

			BeforeEach(func() {
				handlerFunc := func(context.Context, *Command) error {
					return nil
				}

				err = commandBus.Register(handlerFunc)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when register multiple handlerFunc with same command", func() {
			var err error

			BeforeEach(func() {
				for i := 1; i <= 2; i++ {
					handlerFunc := func(context.Context, *Command) error {
						return nil
					}

					err = commandBus.Register(handlerFunc)
				}
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(commandbus.ErrCommandIsAlreadyRegistered))
			})
		})
	})

	Describe("Execute", func() {
		Context("when handler has not registered yet", func() {
			Context("when middleware has not registered yet", func() {
				DescribeTable("when command is invalid",
					func(command interface{}, expectedErr error) {
						err := commandBus.Execute(ctx, command)
						Expect(err).Should(MatchError(expectedErr))
					},
					Entry("when command is nil", nil, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is integer", 5, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is struct", Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is slice of struct", []Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is pointer of pointer of struct", Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
				)

				Context("when command is valid", func() {
					var err error

					BeforeEach(func() {
						command := &Command{}
						err = commandBus.Execute(ctx, command)
					})

					It("should return error", func() {
						Expect(err).Should(MatchError(commandbus.ErrCommandHasNotRegisteredYet))
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

				handlerFunc := func(context.Context, *Command) error {
					handlerExecuted = true
					executedOrders = append(executedOrders, Handler)
					return nil
				}

				err := commandBus.Register(handlerFunc)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Context("when middleware has not registered yet", func() {
				DescribeTable("when command is invalid",
					func(command interface{}, expectedErr error) {
						err := commandBus.Execute(ctx, command)
						Expect(err).Should(MatchError(expectedErr))
					},
					Entry("when command is nil", nil, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is integer", 5, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is struct", Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is slice of struct", []Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is pointer of pointer of struct", Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
				)

				Context("when command is valid", func() {
					var err error

					BeforeEach(func() {
						command := &Command{}
						err = commandBus.Execute(ctx, command)
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
					firstMiddleware := func(commandHandlerFunc commandbus.CommandHandlerFunc[any]) commandbus.CommandHandlerFunc[any] {
						return func(ctx context.Context, command any) error {
							executedOrders = append(executedOrders, FirstMiddlewareBeforeHandler)
							err := commandHandlerFunc(ctx, command)
							executedOrders = append(executedOrders, FirstMiddlewareAfterHandler)
							return err
						}
					}

					secondMiddleware := func(commandHandlerFunc commandbus.CommandHandlerFunc[any]) commandbus.CommandHandlerFunc[any] {
						return func(ctx context.Context, command any) error {
							executedOrders = append(executedOrders, SecondMiddlewareBeforeHandler)
							err := commandHandlerFunc(ctx, command)
							executedOrders = append(executedOrders, SecondMiddlewareAfterHandler)
							return err
						}
					}

					thirdMiddleware := func(commandHandlerFunc commandbus.CommandHandlerFunc[any]) commandbus.CommandHandlerFunc[any] {
						return func(ctx context.Context, command any) error {
							executedOrders = append(executedOrders, ThirdMiddlewareBeforeHandler)
							err := commandHandlerFunc(ctx, command)
							executedOrders = append(executedOrders, ThirdMiddlewareAfterHandler)
							return err
						}
					}

					err := commandBus.Use(firstMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = commandBus.Use(secondMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = commandBus.Use(thirdMiddleware)
					Expect(err).ShouldNot(HaveOccurred())
				})

				DescribeTable("when command is invalid",
					func(command interface{}, expectedErr error) {
						err := commandBus.Execute(ctx, command)
						Expect(err).Should(MatchError(expectedErr))
					},
					Entry("when command is nil", nil, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is integer", 5, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is struct", Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is slice of struct", []Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is pointer of pointer of struct", Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
				)

				Context("when command is valid", func() {
					var err error

					BeforeEach(func() {
						command := &Command{}
						err = commandBus.Execute(ctx, command)
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

				handlerFunc := func(context.Context, *Command) error {
					handlerExecuted = true
					executedOrders = append(executedOrders, Handler)
					return Err
				}

				err := commandBus.Register(handlerFunc)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Context("when middleware has not registered yet", func() {
				DescribeTable("when command is invalid",
					func(command interface{}, expectedErr error) {
						err := commandBus.Execute(ctx, command)
						Expect(err).Should(MatchError(expectedErr))
					},
					Entry("when command is nil", nil, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is integer", 5, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is struct", Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is slice of struct", []Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is pointer of pointer of struct", Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
				)

				Context("when command is valid", func() {
					var err error

					BeforeEach(func() {
						command := &Command{}
						err = commandBus.Execute(ctx, command)
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
					firstMiddleware := func(commandHandlerFunc commandbus.CommandHandlerFunc[any]) commandbus.CommandHandlerFunc[any] {
						return func(ctx context.Context, command any) error {
							executedOrders = append(executedOrders, FirstMiddlewareBeforeHandler)
							err := commandHandlerFunc(ctx, command)
							executedOrders = append(executedOrders, FirstMiddlewareAfterHandler)
							return err
						}
					}

					secondMiddleware := func(commandHandlerFunc commandbus.CommandHandlerFunc[any]) commandbus.CommandHandlerFunc[any] {
						return func(ctx context.Context, command any) error {
							executedOrders = append(executedOrders, SecondMiddlewareBeforeHandler)
							err := commandHandlerFunc(ctx, command)
							executedOrders = append(executedOrders, SecondMiddlewareAfterHandler)
							return err
						}
					}

					thirdMiddleware := func(commandHandlerFunc commandbus.CommandHandlerFunc[any]) commandbus.CommandHandlerFunc[any] {
						return func(ctx context.Context, command any) error {
							executedOrders = append(executedOrders, ThirdMiddlewareBeforeHandler)
							err := commandHandlerFunc(ctx, command)
							executedOrders = append(executedOrders, ThirdMiddlewareAfterHandler)
							return err
						}
					}

					err := commandBus.Use(firstMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = commandBus.Use(secondMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = commandBus.Use(thirdMiddleware)
					Expect(err).ShouldNot(HaveOccurred())
				})

				DescribeTable("when command is invalid",
					func(command interface{}, expectedErr error) {
						err := commandBus.Execute(ctx, command)
						Expect(err).Should(MatchError(expectedErr))
					},
					Entry("when command is nil", nil, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is integer", 5, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is struct", Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is slice of struct", []Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
					Entry("when command is pointer of pointer of struct", Command{}, commandbus.ErrCommandMustBeNonNilPointerOfStruct),
				)

				Context("when command is valid", func() {
					var err error

					BeforeEach(func() {
						command := &Command{}
						err = commandBus.Execute(ctx, command)
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
