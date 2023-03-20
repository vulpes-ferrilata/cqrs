package eventprovider_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEventprovider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Eventprovider Suite")
}
