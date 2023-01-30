package cqrs_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/vulpes-ferrilata/cqrs"
)

type Query struct{}

var _ = Describe("QueryBus", func() {
	var queryBus *cqrs.QueryBus

	BeforeEach(func() {
		queryBus = &cqrs.QueryBus{}
	})

	When("middleware was not registered", func() {
		When("query was not registered", func() {
			It("cannot execute query", func(ctx SpecContext) {
				_, err := queryBus.Execute(ctx, &Query{})
				Expect(err).Should(MatchError(cqrs.ErrHandlerNotFound))
			})
		})

		When("query was registered", func() {
			type Test struct {
				HandlerExecuted bool
			}

			handler := func(ctx context.Context, query *Query) (string, error) {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.HandlerExecuted = true
				}

				return "result", nil
			}

			BeforeEach(func() {
				queryBus.Register(&Query{}, cqrs.WrapQueryHandlerFunc(handler))
			})

			It("can execute query", func(ctx SpecContext) {
				test := &Test{
					HandlerExecuted: false,
				}
				ctxWithTest := context.WithValue(ctx, "test", test)

				result, err := queryBus.Execute(ctxWithTest, &Query{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).Should(BeEquivalentTo("result"))
				Expect(test.HandlerExecuted).Should(BeTrue())
			})
		})
	})

	When("middleware was registered", func() {
		type Test struct {
			HandlerExecuted bool
			Order           []int
		}

		middlewareFunc1 := func(queryHandlerFunc cqrs.QueryHandlerFunc[any, any]) cqrs.QueryHandlerFunc[any, any] {
			return func(ctx context.Context, query any) (interface{}, error) {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.Order = append(test.Order, 1)
				}

				return queryHandlerFunc(ctx, query)
			}
		}

		middlewareFunc2 := func(queryHandlerFunc cqrs.QueryHandlerFunc[any, any]) cqrs.QueryHandlerFunc[any, any] {
			return func(ctx context.Context, query any) (interface{}, error) {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.Order = append(test.Order, 2)
				}

				return queryHandlerFunc(ctx, query)
			}
		}

		middlewareFunc3 := func(queryHandlerFunc cqrs.QueryHandlerFunc[any, any]) cqrs.QueryHandlerFunc[any, any] {
			return func(ctx context.Context, query any) (interface{}, error) {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.Order = append(test.Order, 3)
				}

				return queryHandlerFunc(ctx, query)
			}
		}

		BeforeEach(func() {
			queryBus.Use(
				middlewareFunc1,
				middlewareFunc2,
				middlewareFunc3,
			)
		})

		When("query was not registered", func() {
			It("cannot execute query", func(ctx SpecContext) {
				_, err := queryBus.Execute(ctx, &Query{})
				Expect(err).Should(MatchError(cqrs.ErrHandlerNotFound))
			})
		})

		When("query was registered", func() {
			handler := func(ctx context.Context, query *Query) (string, error) {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.HandlerExecuted = true
					test.Order = append(test.Order, 4)
				}

				return "result", nil
			}

			BeforeEach(func() {
				queryBus.Register(&Query{}, cqrs.WrapQueryHandlerFunc(handler))
			})

			It("can execute query", func(ctx SpecContext) {
				test := &Test{
					HandlerExecuted: false,
					Order:           make([]int, 0),
				}
				ctxWithTest := context.WithValue(ctx, "test", test)

				result, err := queryBus.Execute(ctxWithTest, &Query{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(test.HandlerExecuted).Should(BeTrue())
				Expect(result).Should(BeEquivalentTo("result"))
				Expect(test.Order).Should(BeEquivalentTo([]int{1, 2, 3, 4}))
			})
		})
	})
})
