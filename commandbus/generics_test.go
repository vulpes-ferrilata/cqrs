package commandbus_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/vulpes-ferrilata/cqrs/commandbus"
)

type CommandHandler struct{}

func (c CommandHandler) Handle(ctx context.Context, command *Command) error {
	return nil
}

type CommandBusMock struct {
	mock.Mock
}

func (c CommandBusMock) Use(middlewareFunc commandbus.CommandMiddlewareFunc) error {
	args := c.Called(middlewareFunc)
	return args.Error(0)
}

func (c CommandBusMock) Register(handlerFunc interface{}) error {
	args := c.Called(handlerFunc)
	return args.Error(0)
}

func (c CommandBusMock) Execute(ctx context.Context, command interface{}) error {
	args := c.Called(ctx, command)
	return args.Error(0)
}

var _ = Describe("Generics", func() {
	var ctx context.Context
	var commandBus *CommandBusMock

	BeforeEach(func() {
		ctx = context.Background()
		commandBus = &CommandBusMock{}
	})

	Describe("RegisterCommandHandlerWithCommandBus", func() {
		Context("when CommandBus return error", func() {
			var err error

			expectedErr := errors.New("error")

			BeforeEach(func() {
				commandHandler := &CommandHandler{}

				commandBus.
					On("Register", mock.AnythingOfType("func(context.Context, *commandbus_test.Command) error")).
					Return(expectedErr).
					Once()

				err = commandbus.RegisterCommandHandlerWithCommandBus[Command](commandBus, commandHandler)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(expectedErr))
			})
		})

		Context("when CommandBus not return any error", func() {
			var err error

			BeforeEach(func() {
				commandHandler := &CommandHandler{}

				commandBus.
					On("Register", mock.AnythingOfType("func(context.Context, *commandbus_test.Command) error")).
					Return(nil).
					Once()

				err = commandbus.RegisterCommandHandlerWithCommandBus[Command](commandBus, commandHandler)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("RegisterCommandHandlerFuncWithCommandBus", func() {
		Context("when CommandBus return error", func() {
			var err error

			expectedErr := errors.New("error")

			BeforeEach(func() {
				commandHandlerFunc := func(ctx context.Context, command *Command) error {
					return nil
				}

				commandBus.
					On("Register", mock.AnythingOfType("CommandHandlerFunc[*github.com/vulpes-ferrilata/cqrs/commandbus_test.Command]")).
					Return(expectedErr).
					Once()

				err = commandbus.RegisterCommandHandlerFuncWithCommandBus[Command](commandBus, commandHandlerFunc)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(expectedErr))
			})
		})

		Context("when CommandBus not return any error", func() {
			var err error

			BeforeEach(func() {
				commandHandlerFunc := func(ctx context.Context, command *Command) error {
					return nil
				}

				commandBus.
					On("Register", mock.AnythingOfType("CommandHandlerFunc[*github.com/vulpes-ferrilata/cqrs/commandbus_test.Command]")).
					Return(nil).
					Once()

				err = commandbus.RegisterCommandHandlerFuncWithCommandBus[Command](commandBus, commandHandlerFunc)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("ExecuteCommandWithCommandBus", func() {
		Context("when CommandBus return error", func() {
			var err error

			expectedErr := errors.New("error")

			BeforeEach(func() {
				command := &Command{}

				commandBus.
					On("Execute", ctx, command).
					Return(expectedErr).
					Once()

				err = commandbus.ExecuteCommandWithCommandBus[Command](commandBus, ctx, command)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(expectedErr))
			})
		})

		Context("when CommandBus not return any error", func() {
			var err error

			BeforeEach(func() {
				command := &Command{}

				commandBus.
					On("Execute", ctx, command).
					Return(nil).
					Once()

				err = commandbus.ExecuteCommandWithCommandBus[Command](commandBus, ctx, command)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
