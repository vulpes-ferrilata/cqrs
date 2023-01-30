package cqrs_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/vulpes-ferrilata/cqrs"
)

type Command struct{}

var _ = Describe("CommandBus", func() {
	var commandBus *cqrs.CommandBus

	BeforeEach(func() {
		commandBus = &cqrs.CommandBus{}
	})

	When("middleware was not registered", func() {
		When("command was not registered", func() {
			It("cannot execute command", func(ctx SpecContext) {
				err := commandBus.Execute(ctx, &Command{})
				Expect(err).Should(MatchError(cqrs.ErrHandlerNotFound))
			})
		})

		When("command was registered", func() {
			type Test struct {
				HandlerExecuted bool
			}

			handler := func(ctx context.Context, command *Command) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.HandlerExecuted = true
				}

				return nil
			}

			BeforeEach(func() {
				commandBus.Register(&Command{}, cqrs.WrapCommandHandlerFunc(handler))
			})

			It("can execute command", func(ctx SpecContext) {
				test := &Test{
					HandlerExecuted: false,
				}
				ctxWithTest := context.WithValue(ctx, "test", test)

				err := commandBus.Execute(ctxWithTest, &Command{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(test.HandlerExecuted).Should(BeTrue())
			})
		})
	})

	When("middleware was registered", func() {
		type Test struct {
			HandlerExecuted bool
			Order           []int
		}

		middlewareFunc1 := func(commandHandlerFunc cqrs.CommandHandlerFunc[any]) cqrs.CommandHandlerFunc[any] {
			return func(ctx context.Context, command any) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.Order = append(test.Order, 1)
				}

				return commandHandlerFunc(ctx, command)
			}
		}

		middlewareFunc2 := func(commandHandlerFunc cqrs.CommandHandlerFunc[any]) cqrs.CommandHandlerFunc[any] {
			return func(ctx context.Context, command any) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.Order = append(test.Order, 2)
				}

				return commandHandlerFunc(ctx, command)
			}
		}

		middlewareFunc3 := func(commandHandlerFunc cqrs.CommandHandlerFunc[any]) cqrs.CommandHandlerFunc[any] {
			return func(ctx context.Context, command any) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.Order = append(test.Order, 3)
				}

				return commandHandlerFunc(ctx, command)
			}
		}

		BeforeEach(func() {
			commandBus.Use(
				middlewareFunc1,
				middlewareFunc2,
				middlewareFunc3,
			)
		})

		When("command was not registered", func() {
			It("cannot execute command", func(ctx SpecContext) {
				err := commandBus.Execute(ctx, &Command{})
				Expect(err).Should(MatchError(cqrs.ErrHandlerNotFound))
			})
		})

		When("command was registered", func() {
			handler := func(ctx context.Context, command *Command) error {
				test, ok := ctx.Value("test").(*Test)

				if ok {
					test.HandlerExecuted = true
					test.Order = append(test.Order, 4)
				}

				return nil
			}

			BeforeEach(func() {
				commandBus.Register(&Command{}, cqrs.WrapCommandHandlerFunc(handler))
			})

			It("can execute command", func(ctx SpecContext) {
				test := &Test{
					HandlerExecuted: false,
					Order:           make([]int, 0),
				}
				ctxWithTest := context.WithValue(ctx, "test", test)

				err := commandBus.Execute(ctxWithTest, &Command{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(test.HandlerExecuted).Should(BeTrue())
				Expect(test.Order).Should(BeEquivalentTo([]int{1, 2, 3, 4}))
			})
		})
	})
})
