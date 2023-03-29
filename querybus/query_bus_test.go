package querybus_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vulpes-ferrilata/cqrs/v2/querybus"
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

type Query struct{}

type Result struct{}

var _ = Describe("QueryBus", func() {
	var ctx context.Context
	var queryBus querybus.QueryBus

	BeforeEach(func() {
		ctx = context.Background()
		queryBus = querybus.NewQueryBus()
	})

	Describe("Use", func() {
		Context("when middlewareFunc is nil", func() {
			var err error

			BeforeEach(func() {
				err = queryBus.Use(nil)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(querybus.ErrMiddlewareFuncMustNotBeNil))
			})
		})

		Context("when middlewareFunc is function", func() {
			var err error

			BeforeEach(func() {
				middlewareFunc := func(queryHandlerFunc querybus.QueryHandlerFunc[any, any]) querybus.QueryHandlerFunc[any, any] {
					return func(ctx context.Context, query interface{}) (interface{}, error) {
						return queryHandlerFunc(ctx, query)
					}
				}

				err = queryBus.Use(middlewareFunc)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("Register", func() {
		DescribeTable("when handlerFunc is invalid",
			func(handlerFunc interface{}, expectedErr error) {
				err := queryBus.Register(handlerFunc)
				Expect(err).Should(MatchError(expectedErr))
			},
			Entry("when handlerFunc is nil", nil, querybus.ErrHandlerFuncMustBeNonNilFunction),
			Entry("when handlerFunc is not function", struct{}{}, querybus.ErrHandlerFuncMustBeNonNilFunction),
			Entry("when handlerFunc has 1 argument", func(context.Context) {}, querybus.ErrHandlerFuncMustHaveExactTwoArguments),
			Entry("when handlerFunc has 3 arguments", func(context.Context, *Query, bool) {}, querybus.ErrHandlerFuncMustHaveExactTwoArguments),
			Entry("when first argument of handlerFunc is not context.Context", func(struct{}, *Query) {}, querybus.ErrFirstArgumentOfHandlerMustBeContext),
			Entry("when second argument of handlerFunc is integer", func(context.Context, int) {}, querybus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when second argument of handlerFunc is struct", func(context.Context, Query) {}, querybus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when second argument of handlerFunc is slice of pointer of struct", func(context.Context, []Query) {}, querybus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when second argument of handlerFunc is pointer of pointer of struct", func(context.Context, **Query) {}, querybus.ErrSecondArgumentOfHandlerMustBePointerOfStruct),
			Entry("when handlerFunc has one result", func(context.Context, *Query) {}, querybus.ErrHandlerFuncMustHaveExactTwoResults),
			Entry("when handlerFunc has 3 results", func(context.Context, *Query) (interface{}, interface{}, error) { return nil, nil, nil }, querybus.ErrHandlerFuncMustHaveExactTwoResults),
			Entry("when first result of handlerFunc is integer", func(context.Context, *Query) (int, error) { return 0, nil }, querybus.ErrFirstResultOfHandlerMustBePointerOfStruct),
			Entry("when first result of handlerFunc is struct", func(context.Context, *Query) (Result, error) { return Result{}, nil }, querybus.ErrFirstResultOfHandlerMustBePointerOfStruct),
			Entry("when first result of handlerFunc is slice of pointer of struct", func(context.Context, *Query) ([]*Result, error) { return []*Result{}, nil }, querybus.ErrFirstResultOfHandlerMustBePointerOfStruct),
			Entry("when first result of handlerFunc is pointer of pointer of struct", func(context.Context, *Query) (**Result, error) { return nil, nil }, querybus.ErrFirstResultOfHandlerMustBePointerOfStruct),
			Entry("when second result of handlerFunc is not error", func(context.Context, *Query) (*Result, interface{}) { return nil, nil }, querybus.ErrSecondResultOfHandlerMustBeError),
		)

		Context("when handlerFunc is valid", func() {
			var err error

			BeforeEach(func() {
				handlerFunc := func(context.Context, *Query) (*Result, error) {
					return nil, nil
				}

				err = queryBus.Register(handlerFunc)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when register multiple handlerFunc with same query", func() {
			var err error

			BeforeEach(func() {
				for i := 1; i <= 2; i++ {
					handlerFunc := func(context.Context, *Query) (*Result, error) {
						return nil, nil
					}

					err = queryBus.Register(handlerFunc)
				}
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(querybus.ErrQueryIsAlreadyRegistered))
			})
		})
	})

	Describe("Execute", func() {
		Context("when handler has not registered yet", func() {
			Context("when middleware has not registered yet", func() {
				DescribeTable("when query is invalid",
					func(query interface{}, expectedErr error) {
						result, err := queryBus.Execute(ctx, query)
						Expect(err).Should(MatchError(expectedErr))
						Expect(result).Should(BeNil())
					},
					Entry("when query is nil", nil, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is integer", 5, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is struct", Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is slice of struct", []Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is pointer of pointer of struct", Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
				)

				Context("when query is valid", func() {
					var result interface{}
					var err error

					BeforeEach(func() {
						query := &Query{}
						result, err = queryBus.Execute(ctx, query)
					})

					It("should return error", func() {
						Expect(err).Should(MatchError(querybus.ErrQueryHasNotRegisteredYet))
					})

					It("should return nil result", func() {
						Expect(result).Should(BeNil())
					})
				})
			})
		})

		Context("when handler is registered", func() {
			var handlerExecuted bool
			var executedOrders []HandlerType

			expectedResult := &Result{}

			BeforeEach(func() {
				handlerExecuted = false
				executedOrders = make([]HandlerType, 0)

				handlerFunc := func(context.Context, *Query) (*Result, error) {
					handlerExecuted = true
					executedOrders = append(executedOrders, Handler)
					return expectedResult, nil
				}

				err := queryBus.Register(handlerFunc)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Context("when middleware has not registered yet", func() {
				DescribeTable("when query is invalid",
					func(query interface{}, expectedErr error) {
						result, err := queryBus.Execute(ctx, query)
						Expect(err).Should(MatchError(expectedErr))
						Expect(result).Should(BeNil())
					},
					Entry("when query is nil", nil, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is integer", 5, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is struct", Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is slice of struct", []Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is pointer of pointer of struct", Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
				)

				Context("when query is valid", func() {
					var result interface{}
					var err error

					BeforeEach(func() {
						query := &Query{}
						result, err = queryBus.Execute(ctx, query)
					})

					It("should not return any error", func() {
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should return correct result", func() {
						Expect(result).Should(BeEquivalentTo(expectedResult))
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
					firstMiddleware := func(queryHandlerFunc querybus.QueryHandlerFunc[any, any]) querybus.QueryHandlerFunc[any, any] {
						return func(ctx context.Context, query any) (interface{}, error) {
							executedOrders = append(executedOrders, FirstMiddlewareBeforeHandler)
							result, err := queryHandlerFunc(ctx, query)
							executedOrders = append(executedOrders, FirstMiddlewareAfterHandler)
							return result, err
						}
					}

					secondMiddleware := func(queryHandlerFunc querybus.QueryHandlerFunc[any, any]) querybus.QueryHandlerFunc[any, any] {
						return func(ctx context.Context, query any) (interface{}, error) {
							executedOrders = append(executedOrders, SecondMiddlewareBeforeHandler)
							result, err := queryHandlerFunc(ctx, query)
							executedOrders = append(executedOrders, SecondMiddlewareAfterHandler)
							return result, err
						}
					}

					thirdMiddleware := func(queryHandlerFunc querybus.QueryHandlerFunc[any, any]) querybus.QueryHandlerFunc[any, any] {
						return func(ctx context.Context, query any) (interface{}, error) {
							executedOrders = append(executedOrders, ThirdMiddlewareBeforeHandler)
							result, err := queryHandlerFunc(ctx, query)
							executedOrders = append(executedOrders, ThirdMiddlewareAfterHandler)
							return result, err
						}
					}

					err := queryBus.Use(firstMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = queryBus.Use(secondMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = queryBus.Use(thirdMiddleware)
					Expect(err).ShouldNot(HaveOccurred())
				})

				DescribeTable("when query is invalid",
					func(query interface{}, expectedErr error) {
						result, err := queryBus.Execute(ctx, query)
						Expect(err).Should(MatchError(expectedErr))
						Expect(result).Should(BeNil())
					},
					Entry("when query is nil", nil, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is integer", 5, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is struct", Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is slice of struct", []Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is pointer of pointer of struct", Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
				)

				Context("when query is valid", func() {
					var result interface{}
					var err error

					BeforeEach(func() {
						query := &Query{}
						result, err = queryBus.Execute(ctx, query)
					})

					It("should not return any error", func() {
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should return correct result", func() {
						Expect(result).Should(BeEquivalentTo(expectedResult))
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

			expectedResult := &Result{}
			expectedErr := errors.New("error")

			BeforeEach(func() {
				handlerExecuted = false
				executedOrders = make([]HandlerType, 0)

				handlerFunc := func(context.Context, *Query) (*Result, error) {
					handlerExecuted = true
					executedOrders = append(executedOrders, Handler)
					return expectedResult, expectedErr
				}

				err := queryBus.Register(handlerFunc)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Context("when middleware has not registered yet", func() {
				DescribeTable("when query is invalid",
					func(query interface{}, expectedErr error) {
						result, err := queryBus.Execute(ctx, query)
						Expect(err).Should(MatchError(expectedErr))
						Expect(result).Should(BeNil())
					},
					Entry("when query is nil", nil, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is integer", 5, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is struct", Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is slice of struct", []Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is pointer of pointer of struct", Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
				)

				Context("when query is valid", func() {
					var result interface{}
					var err error

					BeforeEach(func() {
						query := &Query{}
						result, err = queryBus.Execute(ctx, query)
					})

					It("should return error", func() {
						Expect(err).Should(MatchError(expectedErr))
					})

					It("should return correct result", func() {
						Expect(result).Should(BeEquivalentTo(expectedResult))
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
					firstMiddleware := func(queryHandlerFunc querybus.QueryHandlerFunc[any, any]) querybus.QueryHandlerFunc[any, any] {
						return func(ctx context.Context, query any) (interface{}, error) {
							executedOrders = append(executedOrders, FirstMiddlewareBeforeHandler)
							result, err := queryHandlerFunc(ctx, query)
							executedOrders = append(executedOrders, FirstMiddlewareAfterHandler)
							return result, err
						}
					}

					secondMiddleware := func(queryHandlerFunc querybus.QueryHandlerFunc[any, any]) querybus.QueryHandlerFunc[any, any] {
						return func(ctx context.Context, query any) (interface{}, error) {
							executedOrders = append(executedOrders, SecondMiddlewareBeforeHandler)
							result, err := queryHandlerFunc(ctx, query)
							executedOrders = append(executedOrders, SecondMiddlewareAfterHandler)
							return result, err
						}
					}

					thirdMiddleware := func(queryHandlerFunc querybus.QueryHandlerFunc[any, any]) querybus.QueryHandlerFunc[any, any] {
						return func(ctx context.Context, query any) (interface{}, error) {
							executedOrders = append(executedOrders, ThirdMiddlewareBeforeHandler)
							result, err := queryHandlerFunc(ctx, query)
							executedOrders = append(executedOrders, ThirdMiddlewareAfterHandler)
							return result, err
						}
					}

					err := queryBus.Use(firstMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = queryBus.Use(secondMiddleware)
					Expect(err).ShouldNot(HaveOccurred())

					err = queryBus.Use(thirdMiddleware)
					Expect(err).ShouldNot(HaveOccurred())
				})

				DescribeTable("when query is invalid",
					func(query interface{}, expectedErr error) {
						result, err := queryBus.Execute(ctx, query)
						Expect(err).Should(MatchError(expectedErr))
						Expect(result).Should(BeNil())
					},
					Entry("when query is nil", nil, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is integer", 5, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is struct", Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is slice of struct", []Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
					Entry("when query is pointer of pointer of struct", Query{}, querybus.ErrQueryMustBeNonNilPointerOfStruct),
				)

				Context("when query is valid", func() {
					var result interface{}
					var err error

					BeforeEach(func() {
						query := &Query{}
						result, err = queryBus.Execute(ctx, query)
					})

					It("should return error", func() {
						Expect(err).Should(MatchError(expectedErr))
					})

					It("should return correct result", func() {
						Expect(result).Should(BeEquivalentTo(expectedResult))
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
