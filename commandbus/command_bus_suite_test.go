package commandbus_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommandBus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CommandBus Suite")
}
