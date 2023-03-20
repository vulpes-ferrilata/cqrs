package eventprovider_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/vulpes-ferrilata/cqrs/eventprovider"
)

type MockAggregateRoot struct {
	mock.Mock
}

type Event struct{}

func (m MockAggregateRoot) GetEvents() []interface{} {
	args := m.Called()
	return args.Get(0).([]interface{})
}

var _ = Describe("EventProvider", func() {
	var eventProvider *eventprovider.EventProvider
	var aggregateRoot *MockAggregateRoot

	BeforeEach(func() {
		eventProvider = &eventprovider.EventProvider{}
		aggregateRoot = &MockAggregateRoot{}
	})

	DescribeTable("GetEvents",
		func(eventLength int) {
			events := make([]interface{}, 0)
			for i := 1; i <= eventLength; i++ {
				events = append(events, &Event{})
			}
			aggregateRoot.On("GetEvents").Return(events).Once()

			eventProvider.CollectEvents(aggregateRoot)

			Expect(eventProvider.GetEvents()).Should(HaveExactElements(events))
		},
		Entry("with no event", 0),
		Entry("with 1 event", 1),
		Entry("with 2 events", 2),
		Entry("with 50 events", 50),
	)
})
