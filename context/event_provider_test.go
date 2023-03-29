package context_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"context"

	cqrs_context "github.com/vulpes-ferrilata/cqrs/v2/context"
	"github.com/vulpes-ferrilata/cqrs/v2/eventprovider"
)

var _ = Describe("EventProvider", func() {
	var ctx context.Context
	var expectedEventProvider *eventprovider.EventProvider

	BeforeEach(func() {
		ctx = context.Background()
		expectedEventProvider = &eventprovider.EventProvider{}
	})

	Describe("GetEventProvider", func() {
		Context("when context has EventProvider injected", func() {
			var eventProvider *eventprovider.EventProvider
			var err error

			BeforeEach(func() {
				ctx = cqrs_context.WithEventProvider(ctx, expectedEventProvider)
				eventProvider, err = cqrs_context.GetEventProvider(ctx)
			})

			It("should not return any error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should return correct result", func() {
				Expect(eventProvider).Should(BeEquivalentTo(expectedEventProvider))
			})
		})

		Context("when context has not injected EventProvider yet", func() {
			var eventProvider *eventprovider.EventProvider
			var err error

			BeforeEach(func() {
				eventProvider, err = cqrs_context.GetEventProvider(ctx)
			})

			It("should return error", func() {
				Expect(err).Should(MatchError(cqrs_context.ErrEventProviderNotFound))
			})

			It("should return nil result", func() {
				Expect(eventProvider).Should(BeNil())
			})
		})
	})
})
