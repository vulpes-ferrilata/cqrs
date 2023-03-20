package commandbus_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/vulpes-ferrilata/cqrs/commandbus"
)

type CommandHandler struct{}

func (c CommandHandler) Handle(ctx context.Context, command *Command) error {
	return nil
}

var _ = Describe("Generics", func() {
	var ctx context.Context
	var commandBus commandbus.CommandBus

	BeforeEach(func() {
		ctx = context.Background()
		commandBus = commandbus.NewCommandBus()
	})

	Describe("RegisterCommandHandlerWithCommandBus", func() {
		var err error

		BeforeEach(func() {
			commandHandler := &CommandHandler{}
			err = commandbus.RegisterCommandHandlerWithCommandBus[Command](commandBus, commandHandler)
		})

		It("should not return any error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("RegisterCommandHandlerFuncWithCommandBus", func() {
		var err error

		BeforeEach(func() {
			commandHandlerFunc := func(ctx context.Context, command *Command) error {
				return nil
			}
			err = commandbus.RegisterCommandHandlerFuncWithCommandBus[Command](commandBus, commandHandlerFunc)
		})

		It("should not return any error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("ExecuteCommandWithCommandBus", func() {
		var err error

		BeforeEach(func() {
			command := &Command{}
			err = commandbus.ExecuteCommandWithCommandBus[Command](commandBus, ctx, command)
		})

		It("should return error", func() {
			Expect(err).Should(MatchError(commandbus.ErrCommandHasNotRegisteredYet))
		})
	})
})
