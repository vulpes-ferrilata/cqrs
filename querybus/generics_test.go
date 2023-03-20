package querybus_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/vulpes-ferrilata/cqrs/querybus"
)

type QueryHandler struct{}

func (q QueryHandler) Handle(ctx context.Context, query *Query) (*Result, error) {
	return nil, nil
}

var _ = Describe("Generics", func() {
	var ctx context.Context
	var queryBus querybus.QueryBus

	BeforeEach(func() {
		ctx = context.Background()
		queryBus = querybus.NewQueryBus()
	})

	Describe("RegisterQueryHandlerWithQueryBus", func() {
		var err error

		BeforeEach(func() {
			queryHandler := &QueryHandler{}
			err = querybus.RegisterQueryHandlerWithQueryBus[Query, Result](queryBus, queryHandler)
		})

		It("should not return any error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("RegisterQueryHandlerFuncWithQueryBus", func() {
		var err error

		BeforeEach(func() {
			queryHandlerFunc := func(ctx context.Context, query *Query) (*Result, error) {
				return nil, nil
			}
			err = querybus.RegisterQueryHandlerFuncWithQueryBus[Query, Result](queryBus, queryHandlerFunc)
		})

		It("should not return any error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("ExecuteQueryWithQueryBus", func() {
		When("handler has not registered yet", func() {
			var result *Result
			var err error

			BeforeEach(func() {
				query := &Query{}
				result, err = querybus.ExecuteQueryWithQueryBus[Query, Result](queryBus, ctx, query)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(querybus.ErrQueryHasNotRegisteredYet))
			})

			It("should return nil result", func() {
				Expect(result).Should(BeNil())
			})
		})

		When("handler is registered", func() {
			expectedResult := &Result{}

			var result *Result
			var err error

			BeforeEach(func() {
				queryHandlerFunc := func(ctx context.Context, query *Query) (*Result, error) {
					return expectedResult, nil
				}
				err = querybus.RegisterQueryHandlerFuncWithQueryBus[Query, Result](queryBus, queryHandlerFunc)
				Expect(err).ShouldNot(HaveOccurred())

				query := &Query{}
				result, err = querybus.ExecuteQueryWithQueryBus[Query, Result](queryBus, ctx, query)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should return correct result", func() {
				Expect(result).Should(BeEquivalentTo(expectedResult))
			})
		})
	})
})
