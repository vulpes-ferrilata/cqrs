package querybus_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/vulpes-ferrilata/cqrs/v2/querybus"
)

type QueryHandler struct{}

func (q QueryHandler) Handle(ctx context.Context, query *Query) (*Result, error) {
	return nil, nil
}

type QueryBusMock struct {
	mock.Mock
}

func (q QueryBusMock) Use(middlewareFunc querybus.QueryMiddlewareFunc) error {
	args := q.Called(middlewareFunc)
	return args.Error(0)
}

func (q QueryBusMock) Register(handlerFunc interface{}) error {
	args := q.Called(handlerFunc)
	return args.Error(0)
}

func (q QueryBusMock) Execute(ctx context.Context, query interface{}) (interface{}, error) {
	args := q.Called(ctx, query)
	return args.Get(0), args.Error(1)
}

var _ = Describe("Generics", func() {
	var ctx context.Context
	var queryBus *QueryBusMock

	BeforeEach(func() {
		ctx = context.Background()
		queryBus = &QueryBusMock{}
	})

	Describe("RegisterQueryHandlerWithQueryBus", func() {
		Context("when QueryBus return error", func() {
			var err error

			expectedErr := errors.New("error")

			BeforeEach(func() {
				queryHandler := &QueryHandler{}

				queryBus.
					On("Register", mock.AnythingOfType("func(context.Context, *querybus_test.Query) (*querybus_test.Result, error)")).
					Return(expectedErr).
					Once()

				err = querybus.RegisterQueryHandlerWithQueryBus[Query, Result](queryBus, queryHandler)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(expectedErr))
			})
		})

		Context("when QueryBus not return any error", func() {
			var err error

			BeforeEach(func() {
				queryHandler := &QueryHandler{}

				queryBus.
					On("Register", mock.AnythingOfType("func(context.Context, *querybus_test.Query) (*querybus_test.Result, error)")).
					Return(nil).
					Once()

				err = querybus.RegisterQueryHandlerWithQueryBus[Query, Result](queryBus, queryHandler)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("RegisterQueryHandlerFuncWithQueryBus", func() {
		Context("when QueryBus return error", func() {
			var err error

			expectedErr := errors.New("error")

			BeforeEach(func() {
				queryHandlerFunc := func(ctx context.Context, query *Query) (*Result, error) {
					return nil, nil
				}

				queryBus.
					On("Register", mock.AnythingOfType("QueryHandlerFunc[*github.com/vulpes-ferrilata/cqrs/v2/querybus_test.Query,*github.com/vulpes-ferrilata/cqrs/v2/querybus_test.Result]")).
					Return(expectedErr).
					Once()

				err = querybus.RegisterQueryHandlerFuncWithQueryBus[Query, Result](queryBus, queryHandlerFunc)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(expectedErr))
			})
		})

		Context("when QueryBus not return any error", func() {
			var err error

			BeforeEach(func() {
				queryHandlerFunc := func(ctx context.Context, query *Query) (*Result, error) {
					return nil, nil
				}

				queryBus.
					On("Register", mock.AnythingOfType("QueryHandlerFunc[*github.com/vulpes-ferrilata/cqrs/v2/querybus_test.Query,*github.com/vulpes-ferrilata/cqrs/v2/querybus_test.Result]")).
					Return(nil).
					Once()

				err = querybus.RegisterQueryHandlerFuncWithQueryBus[Query, Result](queryBus, queryHandlerFunc)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("ExecuteQueryWithQueryBus", func() {
		Context("when QueryBus return error", func() {
			var result *Result
			var err error

			expectedErr := errors.New("error")

			BeforeEach(func() {
				query := &Query{}

				queryBus.
					On("Execute", ctx, query).
					Return(nil, expectedErr).
					Once()

				result, err = querybus.ExecuteQueryWithQueryBus[Query, Result](queryBus, ctx, query)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(expectedErr))
			})

			It("should return nil result", func() {
				Expect(result).Should(BeNil())
			})
		})

		Context("when QueryBus return result", func() {
			var result *Result
			var err error

			expectedResult := &Result{}

			BeforeEach(func() {
				query := &Query{}

				queryBus.
					On("Execute", ctx, query).
					Return(expectedResult, nil).
					Once()

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
